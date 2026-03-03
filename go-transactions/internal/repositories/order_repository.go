package repositories

import (
	"errors"
	"math"
	"time"

	"github.com/iqmalr-pedia/go-transactions/internal/dto/request"
	"github.com/iqmalr-pedia/go-transactions/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id uint) (*models.Order, error)
	GetByOrderNumber(orderNumber string) (*models.Order, error)
	Update(order *models.Order) error

	ListByUserID(userID uint, query request.ListOrdersQuery) ([]models.Order, int64, error)
	ListAll(query request.ListOrdersQuery) ([]models.Order, int64, error)

	CreateItem(item *models.OrderItem) error
	GetItemByID(itemID uint) (*models.OrderItem, error)
	UpdateItem(item *models.OrderItem) error
	ListItemsByVendorID(vendorID uint, query request.ListVendorOrderItemsQuery) ([]models.OrderItem, int64, error)
	GetItemWithOrder(itemID uint) (*models.OrderItem, *models.Order, error)

	CreatePayment(payment *models.OrderPayment) error
	CreateStatusHistory(history *models.OrderStatusHistory) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.
		Preload("Items").
		Preload("Payments").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.
		Preload("Items").
		Preload("Payments").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("order_number = ?", orderNumber).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) ListByUserID(userID uint, query request.ListOrdersQuery) ([]models.Order, int64, error) {
	return r.listOrders(query, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	})
}

func (r *orderRepository) ListAll(query request.ListOrdersQuery) ([]models.Order, int64, error) {
	return r.listOrders(query, nil)
}

func (r *orderRepository) listOrders(query request.ListOrdersQuery, scope func(*gorm.DB) *gorm.DB) ([]models.Order, int64, error) {
	page, limit := normalizePagination(query.Page, query.Limit)

	db := r.db.Model(&models.Order{})
	if scope != nil {
		db = scope(db)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.PaymentStatus != "" {
		db = db.Where("payment_status = ?", query.PaymentStatus)
	}
	if query.StartDate != "" {
		if t, err := time.Parse("2006-01-02", query.StartDate); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if query.EndDate != "" {
		if t, err := time.Parse("2006-01-02", query.EndDate); err == nil {
			db = db.Where("created_at < ?", t.Add(24*time.Hour))
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortCol := "created_at"
	allowedSorts := map[string]bool{
		"created_at": true, "total_amount": true, "status": true, "order_number": true,
	}
	if query.Sort != "" && allowedSorts[query.Sort] {
		sortCol = query.Sort
	}

	sortOrder := "desc"
	if query.Order == "asc" {
		sortOrder = "asc"
	}

	var orders []models.Order
	err := db.
		Order(sortCol + " " + sortOrder).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) CreateItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *orderRepository) GetItemByID(itemID uint) (*models.OrderItem, error) {
	var item models.OrderItem
	err := r.db.Where("id = ?", itemID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *orderRepository) UpdateItem(item *models.OrderItem) error {
	return r.db.Save(item).Error
}

func (r *orderRepository) GetItemWithOrder(itemID uint) (*models.OrderItem, *models.Order, error) {
	var item models.OrderItem
	err := r.db.Where("id = ?", itemID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var order models.Order
	err = r.db.Where("id = ?", item.OrderID).First(&order).Error
	if err != nil {
		return nil, nil, err
	}

	return &item, &order, nil
}

func (r *orderRepository) ListItemsByVendorID(vendorID uint, query request.ListVendorOrderItemsQuery) ([]models.OrderItem, int64, error) {
	page, limit := normalizePagination(query.Page, query.Limit)

	db := r.db.Model(&models.OrderItem{}).Where("vendor_id = ?", vendorID)

	if query.FulfillmentStatus != "" {
		db = db.Where("fulfillment_status = ?", query.FulfillmentStatus)
	}
	if query.StartDate != "" {
		if t, err := time.Parse("2006-01-02", query.StartDate); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if query.EndDate != "" {
		if t, err := time.Parse("2006-01-02", query.EndDate); err == nil {
			db = db.Where("created_at < ?", t.Add(24*time.Hour))
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortCol := "created_at"
	allowedSorts := map[string]bool{
		"created_at": true, "subtotal": true, "fulfillment_status": true,
	}
	if query.Sort != "" && allowedSorts[query.Sort] {
		sortCol = query.Sort
	}

	sortOrder := "desc"
	if query.Order == "asc" {
		sortOrder = "asc"
	}

	var items []models.OrderItem
	err := db.
		Order(sortCol + " " + sortOrder).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&items).Error

	return items, total, err
}

func (r *orderRepository) CreatePayment(payment *models.OrderPayment) error {
	return r.db.Create(payment).Error
}

func (r *orderRepository) CreateStatusHistory(history *models.OrderStatusHistory) error {
	return r.db.Create(history).Error
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}

func CalcTotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
