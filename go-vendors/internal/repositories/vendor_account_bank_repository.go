package repositories

import (
	"github.com/iqmalr-pedia/go-vendors/internal/models"
	"gorm.io/gorm"
)

type VendorAccountBankRepositoryInterface interface {
	CreateAccountBank(vendor *models.VendorBankAccount) error
	UpdateAccountBank(id uint, val *models.VendorBankAccount) error
	DeleteAccountBank(id, vendor_id uint) error
	GetAccountBankByID(id uint) (*models.VendorBankAccount, error)
	GetAccountBankByVendorID(id uint) (*[]models.VendorBankAccount, error)
}

type VendorAccountBankRepository struct {
	db *gorm.DB
}

func NewVendorAccountBankRepository(db *gorm.DB) VendorAccountBankRepositoryInterface {
	return &VendorAccountBankRepository{db: db}
}

func (a *VendorAccountBankRepository) CreateAccountBank(vendor *models.VendorBankAccount) error {
	return a.db.Create(vendor).Error
}

func (a *VendorAccountBankRepository) GetAccountBankByID(id uint) (*models.VendorBankAccount, error) {
	var accountBank models.VendorBankAccount
	err := a.db.First(&accountBank, id).Error
	return &accountBank, err
}

func (a *VendorAccountBankRepository) GetAccountBankByVendorID(id uint) (*[]models.VendorBankAccount, error) {
	var accountBank []models.VendorBankAccount
	err := a.db.Where("vendor_id =?", id).Find(&accountBank).Error
	return &accountBank, err
}

func (a *VendorAccountBankRepository) UpdateAccountBank(id uint, val *models.VendorBankAccount) error {
	return a.db.Where("id =?", id).Updates(val).Error
}

func (a *VendorAccountBankRepository) DeleteAccountBank(id, vendorId uint) error {
	return a.db.Where("vendor_id =? AND id =?", id, vendorId).Delete(&models.VendorBankAccount{}).Error
}
