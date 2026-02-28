package services

import (
	"fmt"
	"time"

	"github.com/iqmalr-pedia/go-transactions/internal/clients"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/request"
	"github.com/iqmalr-pedia/go-transactions/internal/dto/response"
	"github.com/iqmalr-pedia/go-transactions/internal/models"
	"github.com/iqmalr-pedia/go-transactions/internal/repositories"
	"github.com/iqmalr-pedia/go-transactions/pkg/database"
	"gorm.io/gorm"
)

type CartService interface {
	GetCart(userID *uint, sessionID string) (*response.CartResponse, error)
	AddItem(userID *uint, sessionID string, req *request.AddCartItemRequest) (*response.CartItemResponse, error)
	UpdateItem(userID *uint, sessionID string, itemID uint, req *request.UpdateCartItemRequest) (*response.CartItemResponse, error)
	RemoveItem(userID *uint, sessionID string, itemID uint) error
	ClearCart(userID *uint, sessionID string) error
	MergeCart(userID uint, req *request.MergeCartRequest) (*response.MergeCartResponse, error)
	ValidateCart(userID *uint, sessionID string) (*response.CartValidationResponse, error)
}

type cartService struct {
	cartRepo     repositories.CartRepository
	cartItemRepo repositories.CartItemRepository
	productClient clients.ProductClient
}

func NewCartService(
	cartRepo repositories.CartRepository,
	cartItemRepo repositories.CartItemRepository,
	productClient clients.ProductClient,
) CartService {
	return &cartService{
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productClient: productClient,
	}
}

func (s *cartService) GetCart(userID *uint, sessionID string) (*response.CartResponse, error) {
	cart, err := s.getOrCreateCart(userID, sessionID)
	if err != nil {
		return nil, err
	}

	return s.toCartResponse(cart)
}

func (s *cartService) AddItem(userID *uint, sessionID string, req *request.AddCartItemRequest) (*response.CartItemResponse, error) {
	product, err := s.productClient.GetProductByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if !product.IsPublished {
		return nil, fmt.Errorf("product is not available")
	}

	price := product.Price
	stock := product.Stock
	trackInventory := product.TrackInventory
	allowBackorder := product.AllowBackorder

	if req.VariantID != nil {
		variant, err := s.productClient.GetVariantByID(req.ProductID, *req.VariantID)
		if err != nil {
			return nil, fmt.Errorf("variant not found: %w", err)
		}
		if !variant.IsActive {
			return nil, fmt.Errorf("variant is not active")
		}
		price = variant.Price
		stock = variant.Stock
	}

	if trackInventory && !allowBackorder && req.Quantity > stock {
		return nil, &InsufficientStockError{AvailableStock: stock}
	}

	cart, err := s.getOrCreateCart(userID, sessionID)
	if err != nil {
		return nil, err
	}

	existing, err := s.cartItemRepo.GetByCartAndProduct(cart.ID, req.ProductID, req.VariantID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		newQuantity := existing.Quantity + req.Quantity
		if trackInventory && !allowBackorder && newQuantity > stock {
			return nil, &InsufficientStockError{AvailableStock: stock}
		}

		existing.Quantity = newQuantity
		existing.Price = price
		if err := s.cartItemRepo.Update(existing); err != nil {
			return nil, err
		}

		return s.toCartItemResponse(existing, product, req.VariantID)
	}

	item := &models.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		VariantID: req.VariantID,
		Quantity:  req.Quantity,
		Price:     price,
	}

	if err := s.cartItemRepo.Create(item); err != nil {
		return nil, err
	}

	return s.toCartItemResponse(item, product, req.VariantID)
}

func (s *cartService) UpdateItem(userID *uint, sessionID string, itemID uint, req *request.UpdateCartItemRequest) (*response.CartItemResponse, error) {
	cart, err := s.findCart(userID, sessionID)
	if err != nil {
		return nil, err
	}

	item, err := s.cartItemRepo.GetByID(itemID)
	if err != nil {
		return nil, err
	}

	if item.CartID != cart.ID {
		return nil, fmt.Errorf("cart item does not belong to this cart")
	}

	product, err := s.productClient.GetProductByID(item.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	stock := product.Stock
	trackInventory := product.TrackInventory
	allowBackorder := product.AllowBackorder

	if item.VariantID != nil {
		variant, err := s.productClient.GetVariantByID(item.ProductID, *item.VariantID)
		if err != nil {
			return nil, fmt.Errorf("variant not found: %w", err)
		}
		stock = variant.Stock
	}

	if trackInventory && !allowBackorder && req.Quantity > stock {
		return nil, &InsufficientStockError{AvailableStock: stock}
	}

	item.Quantity = req.Quantity
	item.Price = product.Price
	if item.VariantID != nil {
		variant, _ := s.productClient.GetVariantByID(item.ProductID, *item.VariantID)
		if variant != nil {
			item.Price = variant.Price
		}
	}

	if err := s.cartItemRepo.Update(item); err != nil {
		return nil, err
	}

	return s.toCartItemResponse(item, product, item.VariantID)
}

func (s *cartService) RemoveItem(userID *uint, sessionID string, itemID uint) error {
	cart, err := s.findCart(userID, sessionID)
	if err != nil {
		return err
	}

	item, err := s.cartItemRepo.GetByID(itemID)
	if err != nil {
		return err
	}

	if item.CartID != cart.ID {
		return fmt.Errorf("cart item does not belong to this cart")
	}

	return s.cartItemRepo.Delete(itemID)
}

func (s *cartService) ClearCart(userID *uint, sessionID string) error {
	cart, err := s.findCart(userID, sessionID)
	if err != nil {
		return err
	}

	return s.cartItemRepo.DeleteByCartID(cart.ID)
}

func (s *cartService) MergeCart(userID uint, req *request.MergeCartRequest) (*response.MergeCartResponse, error) {
	guestCart, err := s.cartRepo.GetActiveBySessionID(req.SessionID)
	if err != nil {
		return nil, err
	}
	if guestCart == nil {
		return nil, fmt.Errorf("guest cart not found")
	}

	userCart, err := s.cartRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	if userCart == nil {
		guestCart.UserID = &userID
		guestCart.SessionID = ""
		if err := s.cartRepo.Update(guestCart); err != nil {
			return nil, err
		}

		resp, err := s.toCartResponse(guestCart)
		if err != nil {
			return nil, err
		}
		return &response.MergeCartResponse{
			Message: "Cart merged successfully",
			Cart:    *resp,
		}, nil
	}

	err = database.WithTransaction(database.GetDB(), func(tx *gorm.DB) error {
		for _, guestItem := range guestCart.Items {
			existing, err := s.cartItemRepo.GetByCartAndProduct(userCart.ID, guestItem.ProductID, guestItem.VariantID)
			if err != nil {
				return err
			}

			if existing != nil {
				existing.Quantity += guestItem.Quantity
				if err := s.cartItemRepo.Update(existing); err != nil {
					return err
				}
			} else {
				newItem := &models.CartItem{
					CartID:    userCart.ID,
					ProductID: guestItem.ProductID,
					VariantID: guestItem.VariantID,
					Quantity:  guestItem.Quantity,
					Price:     guestItem.Price,
				}
				if err := s.cartItemRepo.Create(newItem); err != nil {
					return err
				}
			}
		}

		guestCart.Status = "merged"
		return s.cartRepo.Update(guestCart)
	})

	if err != nil {
		return nil, err
	}

	updatedCart, err := s.cartRepo.GetByID(userCart.ID)
	if err != nil {
		return nil, err
	}

	resp, err := s.toCartResponse(updatedCart)
	if err != nil {
		return nil, err
	}
	return &response.MergeCartResponse{
		Message: "Cart merged successfully",
		Cart:    *resp,
	}, nil
}

func (s *cartService) ValidateCart(userID *uint, sessionID string) (*response.CartValidationResponse, error) {
	cart, err := s.findCart(userID, sessionID)
	if err != nil {
		return nil, err
	}

	items, err := s.cartItemRepo.GetByCartID(cart.ID)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	valid := true
	var validationItems []response.CartValidationItemResponse

	for _, item := range items {
		product, err := s.productClient.GetProductByID(item.ProductID)
		if err != nil {
			validationItems = append(validationItems, response.CartValidationItemResponse{
				CartItemID:        item.ID,
				ProductID:         item.ProductID,
				Name:              "Unknown",
				RequestedQuantity: item.Quantity,
				AvailableStock:    0,
				Status:            "unavailable",
			})
			valid = false
			continue
		}

		stock := product.Stock
		if item.VariantID != nil {
			variant, err := s.productClient.GetVariantByID(item.ProductID, *item.VariantID)
			if err == nil {
				stock = variant.Stock
			}
		}

		status := "available"
		if product.TrackInventory && !product.AllowBackorder {
			if stock == 0 {
				status = "out_of_stock"
				valid = false
			} else if item.Quantity > stock {
				status = "insufficient"
				valid = false
			}
		}

		validationItems = append(validationItems, response.CartValidationItemResponse{
			CartItemID:        item.ID,
			ProductID:         item.ProductID,
			Name:              product.Name,
			RequestedQuantity: item.Quantity,
			AvailableStock:    stock,
			Status:            status,
		})
	}

	return &response.CartValidationResponse{
		Valid: valid,
		Items: validationItems,
	}, nil
}

func (s *cartService) getOrCreateCart(userID *uint, sessionID string) (*models.Cart, error) {
	var cart *models.Cart
	var err error

	if userID != nil {
		cart, err = s.cartRepo.GetActiveByUserID(*userID)
	} else if sessionID != "" {
		cart, err = s.cartRepo.GetActiveBySessionID(sessionID)
	} else {
		return nil, fmt.Errorf("user_id or session_id is required")
	}

	if err != nil {
		return nil, err
	}

	if cart != nil {
		return cart, nil
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	cart = &models.Cart{
		UserID:    userID,
		SessionID: sessionID,
		Status:    "active",
		ExpiresAt: &expiresAt,
	}

	if err := s.cartRepo.Create(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) findCart(userID *uint, sessionID string) (*models.Cart, error) {
	var cart *models.Cart
	var err error

	if userID != nil {
		cart, err = s.cartRepo.GetActiveByUserID(*userID)
	} else if sessionID != "" {
		cart, err = s.cartRepo.GetActiveBySessionID(sessionID)
	} else {
		return nil, fmt.Errorf("user_id or session_id is required")
	}

	if err != nil {
		return nil, err
	}

	if cart == nil {
		return nil, fmt.Errorf("cart not found")
	}

	return cart, nil
}

func (s *cartService) toCartResponse(cart *models.Cart) (*response.CartResponse, error) {
	var itemResponses []response.CartItemResponse

	for _, item := range cart.Items {
		product, _ := s.productClient.GetProductByID(item.ProductID)
		itemResp, err := s.toCartItemResponse(&item, product, item.VariantID)
		if err != nil {
			continue
		}
		itemResponses = append(itemResponses, *itemResp)
	}

	if itemResponses == nil {
		itemResponses = []response.CartItemResponse{}
	}

	return &response.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		SessionID: cart.SessionID,
		ExpiresAt: cart.ExpiresAt,
		Items:     itemResponses,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}, nil
}

func (s *cartService) toCartItemResponse(item *models.CartItem, product *clients.ProductInfo, variantID *uint) (*response.CartItemResponse, error) {
	resp := &response.CartItemResponse{
		ID:        item.ID,
		CartID:    item.CartID,
		ProductID: item.ProductID,
		VariantID: item.VariantID,
		Quantity:  item.Quantity,
		Price:     item.Price,
		Status:    "available",
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	if product != nil {
		productResp := &response.CartProductResponse{
			ID:    product.ID,
			Name:  product.Name,
			Slug:  product.Slug,
			Price: product.Price,
			Stock: product.Stock,
		}
		if product.PrimaryImage != nil {
			productResp.PrimaryImage = &response.CartImageResponse{
				ID:       product.PrimaryImage.ID,
				ImageURL: product.PrimaryImage.ImageURL,
			}
		}
		resp.Product = productResp

		if product.TrackInventory && !product.AllowBackorder {
			stock := product.Stock
			if variantID != nil {
				variant, err := s.productClient.GetVariantByID(product.ID, *variantID)
				if err == nil {
					stock = variant.Stock
					resp.Variant = &response.CartVariantResponse{
						ID:    variant.ID,
						Name:  variant.Name,
						Price: variant.Price,
						Stock: variant.Stock,
					}
				}
			}

			if stock == 0 {
				resp.Status = "out_of_stock"
			} else if item.Quantity > stock {
				resp.Status = "insufficient_stock"
			}
		}
	} else {
		resp.Status = "unavailable"
	}

	if variantID != nil && resp.Variant == nil {
		variant, err := s.productClient.GetVariantByID(item.ProductID, *variantID)
		if err == nil {
			resp.Variant = &response.CartVariantResponse{
				ID:    variant.ID,
				Name:  variant.Name,
				Price: variant.Price,
				Stock: variant.Stock,
			}
		}
	}

	return resp, nil
}

type InsufficientStockError struct {
	AvailableStock int
}

func (e *InsufficientStockError) Error() string {
	return fmt.Sprintf("insufficient stock, available: %d", e.AvailableStock)
}
