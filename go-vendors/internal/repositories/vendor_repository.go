package repositories

import (
	"time"

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
	AddUserToVendor(vendorID, userID, invitedBy uint, role string) (*models.VendorUser, error)
	IsUserPartOfVendor(vendorID, userID uint) (bool, error)
	ApproveApplication(vendorID, ownerID, approvedBy uint) error
	RejectApplication(vendorID uint) error
	FindUsersByVendorID(vendorID uint) ([]models.VendorUser, error)
	FindVendorUser(vendorID, userID uint) (*models.VendorUser, error)
	UpdateVendorUser(vendorUser *models.VendorUser) error
	RemoveVendorUser(vendorID, userID uint) error
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

func (r *VendorRepository) AddUserToVendor(vendorID, userID uint, invitedBy uint, role string) (*models.VendorUser, error) {
	now := time.Now()
	vendorUser := &models.VendorUser{
		VendorID:  vendorID,
		UserID:    userID,
		Role:      role,
		IsActive:  true,
		InvitedBy: &invitedBy,
		InvitedAt: &now,
		JoinedAt:  &now,
	}

	if err := r.db.Create(vendorUser).Error; err != nil {
		return nil, err
	}

	return vendorUser, nil
}

func (r *VendorRepository) IsUserPartOfVendor(vendorID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.VendorUser{}).Where("vendor_id = ? AND user_id = ? AND is_active = ?", vendorID, userID, true).Count(&count).Error
	return count > 0, err
}

func (r *VendorRepository) ApproveApplication(vendorID, ownerID, approvedBy uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Vendor{}).Where("id = ?", vendorID).Updates(map[string]interface{}{
			"status":      "active",
			"approved_at": time.Now(),
			"approved_by": approvedBy,
		}).Error; err != nil {
			return err
		}

		now := time.Now()
		vendorUser := &models.VendorUser{
			VendorID: vendorID,
			UserID:   ownerID,
			Role:     "owner",
			IsActive: true,
			JoinedAt: &now,
		}
		if err := tx.Create(vendorUser).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *VendorRepository) RejectApplication(vendorID uint) error {
	return r.db.Model(&models.Vendor{}).Where("id = ?", vendorID).Update("status", "rejected").Error
}

func (r *VendorRepository) FindUsersByVendorID(vendorID uint) ([]models.VendorUser, error) {
	var vendorUsers []models.VendorUser
	err := r.db.Preload("User").Where("vendor_id = ?", vendorID).Find(&vendorUsers).Error
	return vendorUsers, err
}

func (r *VendorRepository) FindVendorUser(vendorID, userID uint) (*models.VendorUser, error) {
	var vendorUser models.VendorUser
	err := r.db.Preload("User").Where("vendor_id = ? AND user_id = ?", vendorID, userID).First(&vendorUser).Error
	return &vendorUser, err
}

func (r *VendorRepository) UpdateVendorUser(vendorUser *models.VendorUser) error {
	return r.db.Save(vendorUser).Error
}

func (r *VendorRepository) RemoveVendorUser(vendorID, userID uint) error {
	return r.db.Where("vendor_id = ? AND user_id = ?", vendorID, userID).Delete(&models.VendorUser{}).Error
}
