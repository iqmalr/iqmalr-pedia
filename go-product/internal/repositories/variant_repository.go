package repositories

import (
	"errors"

	"github.com/iqmalr-pedia/go-product/internal/models"
	"gorm.io/gorm"
)

type VariantRepository interface {
	Create(variant *models.ProductVariant) error
	GetByID(id uint) (*models.ProductVariant, error)
	GetByProductID(productID uint) ([]models.ProductVariant, error)
	Update(variant *models.ProductVariant) error
	Delete(id uint) error
	SKUExistsForProduct(productID uint, sku string, excludeID uint) (bool, error)
}

type variantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) VariantRepository {
	return &variantRepository{db: db}
}

func (r *variantRepository) Create(variant *models.ProductVariant) error {
	return r.db.Create(variant).Error
}

func (r *variantRepository) GetByID(id uint) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.Where("id = ?", id).First(&variant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("variant not found")
		}
		return nil, err
	}
	return &variant, nil
}

func (r *variantRepository) GetByProductID(productID uint) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant
	err := r.db.Where("product_id = ?", productID).Order("created_at ASC").Find(&variants).Error
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *variantRepository) Update(variant *models.ProductVariant) error {
	return r.db.Save(variant).Error
}

func (r *variantRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductVariant{}, id).Error
}

func (r *variantRepository) SKUExistsForProduct(productID uint, sku string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProductVariant{}).Where("product_id = ? AND sku = ?", productID, sku)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
