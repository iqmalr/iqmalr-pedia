package repositories

import (
	"github.com/iqmalr-pedia/go-vendors/internal/models"
	"gorm.io/gorm"
)

type VendorRepositoryInterface interface {
	Create(vendor *models.Vendor) error
	FindByID(id uint) (*models.Vendor, error)
	FindBySlug(slug string) (*models.Vendor, error)
	FindByOwnerID(ownerID uint) (*models.Vendor, error)
	Update(vendor *models.Vendor) error
	UpdateStatus(id uint, status string) error
	FindWithPagination(page, limit int, search, status, sort, order string) ([]models.Vendor, int64, error)
	AddUserToVendor(vendorID, userID uint, role string) error
	IsUserPartOfVendor(vendorID, userID uint) (bool, error)
}

type VendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) VendorRepositoryInterface {
	return &VendorRepository{db: db}
}

func (r *VendorRepository) Create(vendor *models.Vendor) error {
	return r.db.Create(vendor).Error
}

func (r *VendorRepository) FindByID(id uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.First(&vendor, id).Error
	return &vendor, err
}

func (r *VendorRepository) FindBySlug(slug string) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Where("slug = ?", slug).First(&vendor).Error
	return &vendor, err
}

func (r *VendorRepository) FindByOwnerID(ownerID uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Where("owner_id = ?", ownerID).First(&vendor).Error
	return &vendor, err
}

func (r *VendorRepository) Update(vendor *models.Vendor) error {
	return r.db.Save(vendor).Error
}

func (r *VendorRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Vendor{}).Where("id = ?", id).Update("status", status).Error
}

func (r *VendorRepository) FindWithPagination(page, limit int, search, status, sort, order string) ([]models.Vendor, int64, error) {
	var vendors []models.Vendor
	var total int64

	query := r.db.Model(&models.Vendor{})

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortClause := sort + " " + order
	query = query.Order(sortClause)

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&vendors).Error; err != nil {
		return nil, 0, err
	}

	return vendors, total, nil
}

func (r *VendorRepository) AddUserToVendor(vendorID, userID uint, role string) error {
	vendorUser := &models.VendorUser{
		VendorID: vendorID,
		UserID:   userID,
		Role:     role,
		IsActive: true,
	}
	return r.db.Create(vendorUser).Error
}

func (r *VendorRepository) IsUserPartOfVendor(vendorID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.VendorUser{}).Where("vendor_id = ? AND user_id = ? AND is_active = ?", vendorID, userID, true).Count(&count).Error
	return count > 0, err
}
