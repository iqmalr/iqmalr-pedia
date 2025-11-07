package repositories

import (
	"errors"

	"github.com/iqmalr-pedia/go-product/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	GetBySlug(slug string) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	List(parentID *uint, level *int, isActive *bool, sort string, order string) ([]models.Category, error)
	GetTree(isActive *bool) ([]models.Category, error)
	HasChildren(id uint) (bool, error)
	HasProducts(id uint) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *categoryRepository) List(parentID *uint, level *int, isActive *bool, sort string, order string) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.Model(&models.Category{})

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	}

	if level != nil {
		query = query.Where("level = ?", *level)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if sort == "" {
		sort = "sort_order"
	}

	if order == "" {
		order = "asc"
	}

	err := query.Order(sort + " " + order).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetTree(isActive *bool) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.Where("parent_id IS NULL")

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Preload("Children.Children.Children").Order("sort_order asc").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *categoryRepository) HasProducts(id uint) (bool, error) {
	// This is a placeholder for when products are implemented
	// For now, we'll just return false
	return false, nil
}
