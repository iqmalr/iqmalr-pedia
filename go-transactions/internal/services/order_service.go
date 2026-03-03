package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/iqmalr-pedia/go-transactions/internal/clients"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/request"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/response"
	"github.com/iqmalr-pedia/go-transactions/internal/models"
	"github.com/iqmalr-pedia/go-transactions/internal/repositories"
	"github.com/iqmalr-pedia/go-transactions/pkg/database"
	"gorm.io/gorm"
)

var (
	validOrderStatuses      = map[string]bool{"pending": true, "processing": true, "shipped": true, "delivered": true, "cancelled": true}
	validPaymentStatuses    = map[string]bool{"pending": true, "paid": true, "failed": true, "refunded": true}
	validFulfillmentStatuses = map[string]bool{"pending": true, "processing": true, "shipped": true, "delivered": true, "cancelled": true}
	validShippingMethods    = map[string]float64{"regular": 15000, "express": 30000, "same_day": 50000, "pickup": 0}
	validPaymentMethods     = map[string]bool{"bank_transfer": true, "e_wallet": true, "cod": true, "credit_card": true}
)

const taxRate = 0.11 // PPN 11%

type OrderService interface {
	CreateOrder(userID uint, authToken string, req *request.CreateOrderRequest) (*response.OrderResponse, error)
	GetOrderByID(userID uint, role string, orderID uint) (*response.OrderResponse, error)
	GetOrderByOrderNumber(userID uint, role string, orderNumber string) (*response.OrderResponse, error)
	ListOrders(userID uint, role string, query request.ListOrdersQuery) (*response.OrderListResponse, error)
	CancelOrder(userID uint, role string, orderID uint, req *request.CancelOrderRequest) error
	UpdateOrderStatus(adminUserID uint, orderID uint, req *request.UpdateOrderStatusRequest) error
	UpdateOrderItemFulfillment(userID uint, role string, orderID uint, itemID uint, req *request.UpdateFulfillmentStatusRequest) error
	GetVendorOrderItems(userID uint, role string, vendorID uint, query request.ListVendorOrderItemsQuery) (*response.VendorOrderItemListResponse, error)
}

type orderService struct {
	orderRepo     repositories.OrderRepository
	cartRepo      repositories.CartRepository
	cartItemRepo  repositories.CartItemRepository
	productClient clients.ProductClient
	authClient    clients.AuthClient
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	cartRepo repositories.CartRepository,
	cartItemRepo repositories.CartItemRepository,
	productClient clients.ProductClient,
	authClient clients.AuthClient,
) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productClient: productClient,
		authClient:    authClient,
	}
}

func (s *orderService) CreateOrder(userID uint, authToken string, req *request.CreateOrderRequest) (*response.OrderResponse, error) {
	shippingCost, ok := validShippingMethods[req.ShippingMethod]
	if !ok {
		return nil, fmt.Errorf("invalid shipping method: must be one of regular, express, same_day, pickup")
	}

	if !validPaymentMethods[req.PaymentMethod] {
		return nil, fmt.Errorf("invalid payment method: must be one of bank_transfer, e_wallet, cod, credit_card")
	}

	var shippingName, shippingPhone, addr1, addr2, city, state, postalCode, country string

	if req.ShippingAddressID != nil && *req.ShippingAddressID > 0 {
		address, err := s.authClient.GetUserAddressByID(authToken, *req.ShippingAddressID)
		if err != nil {
			return nil, fmt.Errorf("failed to get shipping address: %w", err)
		}
		shippingName = address.RecipientName
		shippingPhone = address.Phone
		addr1 = address.AddressLine1
		addr2 = address.AddressLine2
		city = address.City
		state = address.State
		postalCode = address.PostalCode
		country = address.Country
	} else if req.ShippingAddress != nil {
		shippingName = req.ShippingAddress.RecipientName
		shippingPhone = req.ShippingAddress.Phone
		addr1 = req.ShippingAddress.AddressLine1
		addr2 = req.ShippingAddress.AddressLine2
		city = req.ShippingAddress.City
		state = req.ShippingAddress.State
		postalCode = req.ShippingAddress.PostalCode
		country = req.ShippingAddress.Country
		if country == "" {
			country = "Indonesia"
		}
	} else {
		return nil, fmt.Errorf("shipping_address_id or shipping_address is required")
	}

	cart, err := s.cartRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	type orderItemData struct {
		CartItem   models.CartItem
		Product    *clients.ProductInfo
		VariantName string
	}

	var itemsData []orderItemData
	var subtotal float64

	for _, ci := range cart.Items {
		product, err := s.productClient.GetProductByID(ci.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", ci.ProductID, err)
		}
		if !product.IsPublished {
			return nil, fmt.Errorf("product '%s' is not available", product.Name)
		}

		price := product.Price
		stock := product.Stock
		variantName := ""

		if ci.VariantID != nil {
			variant, err := s.productClient.GetVariantByID(ci.ProductID, *ci.VariantID)
			if err != nil {
				return nil, fmt.Errorf("variant %d for product '%s' not found: %w", *ci.VariantID, product.Name, err)
			}
			if !variant.IsActive {
				return nil, fmt.Errorf("variant '%s' for product '%s' is not active", variant.Name, product.Name)
			}
			price = variant.Price
			stock = variant.Stock
			variantName = variant.Name
		}

		if product.TrackInventory && !product.AllowBackorder && ci.Quantity > stock {
			return nil, fmt.Errorf("insufficient stock for product '%s' (requested: %d, available: %d)", product.Name, ci.Quantity, stock)
		}

		itemSubtotal := price * float64(ci.Quantity)
		subtotal += itemSubtotal

		itemsData = append(itemsData, orderItemData{
			CartItem:    ci,
			Product:     product,
			VariantName: variantName,
		})
	}

	taxAmount := math.Round(subtotal*taxRate*100) / 100
	totalAmount := subtotal + shippingCost + taxAmount

	orderNumber := generateOrderNumber()

	order := &models.Order{
		OrderNumber:          orderNumber,
		UserID:               userID,
		ShippingName:         shippingName,
		ShippingPhone:        shippingPhone,
		ShippingAddressLine1: addr1,
		ShippingAddressLine2: addr2,
		ShippingCity:         city,
		ShippingState:        state,
		ShippingPostalCode:   postalCode,
		ShippingCountry:      country,
		Subtotal:             subtotal,
		ShippingCost:         shippingCost,
		TaxAmount:            taxAmount,
		DiscountAmount:       0,
		TotalAmount:          totalAmount,
		PaymentMethod:        req.PaymentMethod,
		PaymentStatus:        "pending",
		Status:               "pending",
		Notes:                req.Notes,
		CouponCode:           req.CouponCode,
		ShippingMethod:       req.ShippingMethod,
	}

	err = database.WithTransaction(database.GetDB(), func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		for _, data := range itemsData {
			vendorID := data.Product.VendorID
			vendorName := ""
			if data.Product.Vendor != nil {
				vendorID = data.Product.Vendor.ID
				vendorName = data.Product.Vendor.Name
			}
			_ = vendorName

			item := models.OrderItem{
				OrderID:           order.ID,
				VendorID:          vendorID,
				ProductID:         data.CartItem.ProductID,
				ProductName:       data.Product.Name,
				VariantID:         data.CartItem.VariantID,
				VariantName:       data.VariantName,
				SKU:               data.Product.SKU,
				Quantity:          data.CartItem.Quantity,
				UnitPrice:         data.CartItem.Price,
				Subtotal:          data.CartItem.Price * float64(data.CartItem.Quantity),
				FulfillmentStatus: "pending",
			}
			if err := tx.Create(&item).Error; err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}

		statusHistory := models.OrderStatusHistory{
			OrderID:   order.ID,
			Status:    "pending",
			Notes:     "Order created",
			CreatedBy: userID,
		}
		if err := tx.Create(&statusHistory).Error; err != nil {
			return fmt.Errorf("failed to create status history: %w", err)
		}

		cart.Status = "ordered"
		if err := tx.Save(cart).Error; err != nil {
			return fmt.Errorf("failed to update cart status: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	created, err := s.orderRepo.GetByID(order.ID)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(created), nil
}

func (s *orderService) GetOrderByID(userID uint, role string, orderID uint) (*response.OrderResponse, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	if role != "admin" && order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) GetOrderByOrderNumber(userID uint, role string, orderNumber string) (*response.OrderResponse, error) {
	order, err := s.orderRepo.GetByOrderNumber(orderNumber)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	if role != "admin" && order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) ListOrders(userID uint, role string, query request.ListOrdersQuery) (*response.OrderListResponse, error) {
	if query.Status != "" && !validOrderStatuses[query.Status] {
		return nil, fmt.Errorf("invalid status filter")
	}
	if query.PaymentStatus != "" && !validPaymentStatuses[query.PaymentStatus] {
		return nil, fmt.Errorf("invalid payment_status filter")
	}

	var orders []models.Order
	var total int64
	var err error

	if role == "admin" {
		orders, total, err = s.orderRepo.ListAll(query)
	} else {
		orders, total, err = s.orderRepo.ListByUserID(userID, query)
	}
	if err != nil {
		return nil, err
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var items []response.OrderListItemResponse
	for _, o := range orders {
		items = append(items, response.OrderListItemResponse{
			ID:             o.ID,
			OrderNumber:    o.OrderNumber,
			UserID:         o.UserID,
			ShippingName:   o.ShippingName,
			ShippingCity:   o.ShippingCity,
			Subtotal:       o.Subtotal,
			ShippingCost:   o.ShippingCost,
			TaxAmount:      o.TaxAmount,
			DiscountAmount: o.DiscountAmount,
			TotalAmount:    o.TotalAmount,
			PaymentMethod:  o.PaymentMethod,
			PaymentStatus:  o.PaymentStatus,
			Status:         o.Status,
			CreatedAt:      o.CreatedAt,
		})
	}

	if items == nil {
		items = []response.OrderListItemResponse{}
	}

	return &response.OrderListResponse{
		Data: items,
		Pagination: response.PaginationResponse{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: repositories.CalcTotalPages(total, limit),
		},
	}, nil
}

func (s *orderService) CancelOrder(userID uint, role string, orderID uint, req *request.CancelOrderRequest) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	if role != "admin" && order.UserID != userID {
		return fmt.Errorf("order not found")
	}

	if order.Status != "pending" && order.Status != "processing" {
		return fmt.Errorf("order cannot be cancelled (current status: %s)", order.Status)
	}

	return database.WithTransaction(database.GetDB(), func(tx *gorm.DB) error {
		order.Status = "cancelled"
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		notes := "Order cancelled"
		if req.Reason != "" {
			notes = "Order cancelled: " + req.Reason
		}

		history := models.OrderStatusHistory{
			OrderID:   order.ID,
			Status:    "cancelled",
			Notes:     notes,
			CreatedBy: userID,
		}
		return tx.Create(&history).Error
	})
}

func (s *orderService) UpdateOrderStatus(adminUserID uint, orderID uint, req *request.UpdateOrderStatusRequest) error {
	if !validOrderStatuses[req.Status] {
		return fmt.Errorf("invalid status: must be one of pending, processing, shipped, delivered, cancelled")
	}

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	return database.WithTransaction(database.GetDB(), func(tx *gorm.DB) error {
		order.Status = req.Status
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		history := models.OrderStatusHistory{
			OrderID:   order.ID,
			Status:    req.Status,
			Notes:     req.Notes,
			CreatedBy: adminUserID,
		}
		return tx.Create(&history).Error
	})
}

func (s *orderService) UpdateOrderItemFulfillment(userID uint, role string, orderID uint, itemID uint, req *request.UpdateFulfillmentStatusRequest) error {
	if !validFulfillmentStatuses[req.FulfillmentStatus] {
		return fmt.Errorf("invalid fulfillment_status: must be one of pending, processing, shipped, delivered, cancelled")
	}

	item, order, err := s.orderRepo.GetItemWithOrder(itemID)
	if err != nil {
		return err
	}
	if item == nil || order == nil {
		return fmt.Errorf("order item not found")
	}
	if order.ID != orderID {
		return fmt.Errorf("order item does not belong to this order")
	}

	if role != "admin" && role != "vendor" {
		return fmt.Errorf("unauthorized: only admin or vendor can update fulfillment status")
	}

	now := time.Now()
	item.FulfillmentStatus = req.FulfillmentStatus
	if req.TrackingNumber != "" {
		item.TrackingNumber = req.TrackingNumber
	}
	if req.FulfillmentStatus == "shipped" {
		item.ShippedAt = &now
	}
	if req.FulfillmentStatus == "delivered" {
		item.DeliveredAt = &now
	}

	return s.orderRepo.UpdateItem(item)
}

func (s *orderService) GetVendorOrderItems(userID uint, role string, vendorID uint, query request.ListVendorOrderItemsQuery) (*response.VendorOrderItemListResponse, error) {
	if role != "admin" && role != "vendor" {
		return nil, fmt.Errorf("unauthorized: only admin or vendor can access vendor order items")
	}

	if query.FulfillmentStatus != "" && !validFulfillmentStatuses[query.FulfillmentStatus] {
		return nil, fmt.Errorf("invalid fulfillment_status filter")
	}

	items, total, err := s.orderRepo.ListItemsByVendorID(vendorID, query)
	if err != nil {
		return nil, err
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orderCache := make(map[uint]*models.Order)
	var respItems []response.VendorOrderItemResponse

	for _, item := range items {
		var orderInfo *response.VendorOrderItemOrderInfo
		if cachedOrder, ok := orderCache[item.OrderID]; ok {
			orderInfo = &response.VendorOrderItemOrderInfo{
				ID:          cachedOrder.ID,
				OrderNumber: cachedOrder.OrderNumber,
				UserID:      cachedOrder.UserID,
				Status:      cachedOrder.Status,
				CreatedAt:   cachedOrder.CreatedAt,
			}
		} else {
			order, err := s.orderRepo.GetByID(item.OrderID)
			if err == nil && order != nil {
				orderCache[item.OrderID] = order
				orderInfo = &response.VendorOrderItemOrderInfo{
					ID:          order.ID,
					OrderNumber: order.OrderNumber,
					UserID:      order.UserID,
					Status:      order.Status,
					CreatedAt:   order.CreatedAt,
				}
			}
		}

		respItems = append(respItems, response.VendorOrderItemResponse{
			ID:                item.ID,
			OrderID:           item.OrderID,
			Order:             orderInfo,
			VendorID:          item.VendorID,
			ProductID:         item.ProductID,
			ProductName:       item.ProductName,
			VariantID:         item.VariantID,
			VariantName:       item.VariantName,
			SKU:               item.SKU,
			Quantity:          item.Quantity,
			UnitPrice:         item.UnitPrice,
			Subtotal:          item.Subtotal,
			FulfillmentStatus: item.FulfillmentStatus,
			TrackingNumber:    item.TrackingNumber,
			ShippedAt:         item.ShippedAt,
			DeliveredAt:       item.DeliveredAt,
			CreatedAt:         item.CreatedAt,
		})
	}

	if respItems == nil {
		respItems = []response.VendorOrderItemResponse{}
	}

	return &response.VendorOrderItemListResponse{
		Data: respItems,
		Pagination: response.PaginationResponse{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: repositories.CalcTotalPages(total, limit),
		},
	}, nil
}

func (s *orderService) toOrderResponse(order *models.Order) *response.OrderResponse {
	var items []response.OrderItemResponse
	for _, item := range order.Items {
		items = append(items, response.OrderItemResponse{
			ID:        item.ID,
			VendorID:  item.VendorID,
			Vendor: &response.OrderVendorResponse{
				ID:   item.VendorID,
				Name: s.getVendorName(item.VendorID),
			},
			ProductID:         item.ProductID,
			ProductName:       item.ProductName,
			VariantID:         item.VariantID,
			VariantName:       item.VariantName,
			SKU:               item.SKU,
			Quantity:          item.Quantity,
			UnitPrice:         item.UnitPrice,
			Subtotal:          item.Subtotal,
			FulfillmentStatus: item.FulfillmentStatus,
			TrackingNumber:    item.TrackingNumber,
			ShippedAt:         item.ShippedAt,
			DeliveredAt:       item.DeliveredAt,
			CreatedAt:         item.CreatedAt,
		})
	}
	if items == nil {
		items = []response.OrderItemResponse{}
	}

	var payments []response.OrderPaymentResponse
	for _, p := range order.Payments {
		payments = append(payments, response.OrderPaymentResponse{
			ID:              p.ID,
			PaymentMethod:   p.PaymentMethod,
			PaymentGateway:  p.PaymentGateway,
			TransactionID:   p.TransactionID,
			Amount:          p.Amount,
			Status:          p.Status,
			PaymentProofURL: p.PaymentProofURL,
			PaidAt:          p.PaidAt,
			CreatedAt:       p.CreatedAt,
		})
	}
	if payments == nil {
		payments = []response.OrderPaymentResponse{}
	}

	var statusHistory []response.OrderStatusHistoryResponse
	for _, sh := range order.StatusHistory {
		statusHistory = append(statusHistory, response.OrderStatusHistoryResponse{
			ID:        sh.ID,
			Status:    sh.Status,
			Notes:     sh.Notes,
			CreatedBy: sh.CreatedBy,
			CreatedAt: sh.CreatedAt,
		})
	}
	if statusHistory == nil {
		statusHistory = []response.OrderStatusHistoryResponse{}
	}

	return &response.OrderResponse{
		ID:                   order.ID,
		OrderNumber:          order.OrderNumber,
		UserID:               order.UserID,
		ShippingName:         order.ShippingName,
		ShippingPhone:        order.ShippingPhone,
		ShippingAddressLine1: order.ShippingAddressLine1,
		ShippingAddressLine2: order.ShippingAddressLine2,
		ShippingCity:         order.ShippingCity,
		ShippingState:        order.ShippingState,
		ShippingPostalCode:   order.ShippingPostalCode,
		ShippingCountry:      order.ShippingCountry,
		Subtotal:             order.Subtotal,
		ShippingCost:         order.ShippingCost,
		TaxAmount:            order.TaxAmount,
		DiscountAmount:       order.DiscountAmount,
		TotalAmount:          order.TotalAmount,
		PaymentMethod:        order.PaymentMethod,
		PaymentStatus:        order.PaymentStatus,
		PaidAt:               order.PaidAt,
		Status:               order.Status,
		Notes:                order.Notes,
		CouponCode:           order.CouponCode,
		Items:                items,
		Payments:             payments,
		StatusHistory:        statusHistory,
		CreatedAt:            order.CreatedAt,
		UpdatedAt:            order.UpdatedAt,
	}
}

func (s *orderService) getVendorName(vendorID uint) string {
	// Best-effort vendor name lookup via product service is not practical here
	// since we only have vendorID but not productID. Return empty for now;
	// vendor name will be populated properly when we add a dedicated vendor client.
	return ""
}

func generateOrderNumber() string {
	now := time.Now()
	datePart := now.Format("20060102")
	suffix := strings.ToUpper(fmt.Sprintf("%05X", now.UnixNano()%0xFFFFF))
	return fmt.Sprintf("ORD-%s-%s", datePart, suffix)
}
