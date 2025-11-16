package repositories

import (
	"errors"
	"time"

	"github.com/iqmalr-pedia/go-product/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	GetByUUID(uuid string) (*models.Product, error)
	GetBySlug(slug string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	List(page, limit int, search string, vendorID *uint, categoryID *uint, minPrice, maxPrice float64, status string, isFeatured, isPublished *bool, sort, order string) ([]models.Product, int64, error)
	UpdateStatus(id uint, status string) error
	Publish(id uint) error
	Unpublish(id uint) error
	CountByVendor(vendorID uint) (int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	// HAPUS .Preload("Vendor") dari sini karena kita tidak lagi memuatnya dari DB
	err := r.db.Preload("Categories").Preload("Images").Preload("Variants").Where("id = ?", id).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByUUID(uuid string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Vendor").Preload("Categories").Preload("Images").Preload("Variants").Where("uuid = ?", uuid).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySlug(slug string) (*models.Product, error) {
	var product models.Product
	// HAPUS .Preload("Vendor") dari sini
	err := r.db.Preload("Categories").Preload("Images").Preload("Variants").Where("slug = ?", slug).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) List(page, limit int, search string, vendorID *uint, categoryID *uint, minPrice, maxPrice float64, status string, isFeatured, isPublished *bool, sort, order string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Preload("Vendor").Preload("Images", "is_primary = ?", true)

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if vendorID != nil {
		query = query.Where("vendor_id = ?", *vendorID)
	}

	if categoryID != nil {
		query = query.Joins("INNER JOIN product_categories ON products.id = product_categories.product_id").
			Where("product_categories.category_id = ?", *categoryID)
	}

	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}

	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if isFeatured != nil {
		query = query.Where("is_featured = ?", *isFeatured)
	}

	if isPublished != nil {
		query = query.Where("is_published = ?", *isPublished)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	offset := (page - 1) * limit
	err = query.Order(sort + " " + order).Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("status", status).Error
}

func (r *productRepository) Publish(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Product{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_published": true,
		"published_at": &now,
	}).Error
}

func (r *productRepository) Unpublish(id uint) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_published": false,
		"published_at": nil,
	}).Error
}

func (r *productRepository) CountByVendor(vendorID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("vendor_id = ?", vendorID).Count(&count).Error
	return count, err
}
