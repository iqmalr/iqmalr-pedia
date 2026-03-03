package repositories

import (
	"errors"

	"github.com/iqmalr-pedia/go-transactions/internal/models"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	Create(item *models.CartItem) error
	GetByID(id uint) (*models.CartItem, error)
	GetByCartID(cartID uint) ([]models.CartItem, error)
	GetByCartAndProduct(cartID, productID uint, variantID *uint) (*models.CartItem, error)
	Update(item *models.CartItem) error
	Delete(id uint) error
	DeleteByCartID(cartID uint) error
}

type cartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) CartItemRepository {
	return &cartItemRepository{db: db}
}

func (r *cartItemRepository) Create(item *models.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartItemRepository) GetByID(id uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart item not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartItemRepository) GetByCartID(cartID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.Where("cart_id = ?", cartID).Order("created_at ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartItemRepository) GetByCartAndProduct(cartID, productID uint, variantID *uint) (*models.CartItem, error) {
	var item models.CartItem
	query := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID)

	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	err := query.First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartItemRepository) Update(item *models.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}

func (r *cartItemRepository) DeleteByCartID(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}
