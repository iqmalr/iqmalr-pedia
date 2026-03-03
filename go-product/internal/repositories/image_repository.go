package repositories

import (
	"errors"

	"github.com/iqmalr-pedia/go-product/internal/models"
	"gorm.io/gorm"
)

type ImageRepository interface {
	Create(image *models.ProductImage) error
	GetByID(id uint) (*models.ProductImage, error)
	GetByProductID(productID uint) ([]models.ProductImage, error)
	Update(image *models.ProductImage) error
	Delete(id uint) error
	ClearPrimaryByProductID(productID uint) error
}

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) Create(image *models.ProductImage) error {
	return r.db.Create(image).Error
}

func (r *imageRepository) GetByID(id uint) (*models.ProductImage, error) {
	var image models.ProductImage
	err := r.db.Where("id = ?", id).First(&image).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("image not found")
		}
		return nil, err
	}
	return &image, nil
}

func (r *imageRepository) GetByProductID(productID uint) ([]models.ProductImage, error) {
	var images []models.ProductImage
	err := r.db.Where("product_id = ?", productID).Order("sort_order ASC, created_at ASC").Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (r *imageRepository) Update(image *models.ProductImage) error {
	return r.db.Save(image).Error
}

func (r *imageRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductImage{}, id).Error
}

func (r *imageRepository) ClearPrimaryByProductID(productID uint) error {
	return r.db.Model(&models.ProductImage{}).
		Where("product_id = ?", productID).
		Update("is_primary", false).Error
}
