package repositories

import (
	"errors"

	"github.com/iqmalr-pedia/go-auth/v2/internal/models"
	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) Create(address *models.UserAddress) error {
	// Jika alamat baru ditandai sebagai default, set alamat default lainnya menjadi false
	if address.IsDefault {
		if err := r.db.Model(&models.UserAddress{}).Where("user_id = ?", address.UserID).Update("is_default", false).Error; err != nil {
			return err
		}
	}
	return r.db.Create(address).Error
}

func (r *AddressRepository) FindByUserID(userID uint) ([]models.UserAddress, error) {
	var addresses []models.UserAddress
	err := r.db.Where("user_id = ?", userID).Find(&addresses).Error
	return addresses, err
}

func (r *AddressRepository) FindByIDAndUserID(id, userID uint) (*models.UserAddress, error) {
	var address models.UserAddress
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ditemukan, bukan error
		}
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) Update(address *models.UserAddress) error {
	if address.IsDefault {
		if err := r.db.Model(&models.UserAddress{}).Where("user_id = ? AND id != ?", address.UserID, address.ID).Update("is_default", false).Error; err != nil {
			return err
		}
	}
	return r.db.Save(address).Error
}

func (r *AddressRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserAddress{}).Error
}

func (r *AddressRepository) SetDefault(id, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.UserAddress{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.UserAddress{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", true).Error; err != nil {
			return err
		}

		return nil
	})
}
