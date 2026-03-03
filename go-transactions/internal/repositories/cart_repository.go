package repositories

import (
	"errors"

	"github.com/iqmalr-pedia/go-transactions/internal/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *models.Cart) error
	GetByID(id uint) (*models.Cart, error)
	GetByUserID(userID uint) (*models.Cart, error)
	GetBySessionID(sessionID string) (*models.Cart, error)
	Update(cart *models.Cart) error
	Delete(id uint) error
	GetActiveByUserID(userID uint) (*models.Cart, error)
	GetActiveBySessionID(sessionID string) (*models.Cart, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) GetByID(id uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items").Where("id = ?", id).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetBySessionID(sessionID string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items").Where("session_id = ?", sessionID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Update(cart *models.Cart) error {
	return r.db.Save(cart).Error
}

func (r *cartRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cart{}, id).Error
}

func (r *cartRepository) GetActiveByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items").Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetActiveBySessionID(sessionID string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items").Where("session_id = ? AND status = ?", sessionID, "active").First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}
