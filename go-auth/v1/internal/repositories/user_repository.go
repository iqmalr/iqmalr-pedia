package repositories

import (
	"time"

	"github.com/iqmalr-pedia/go-auth/v1/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindUserByUUID(uuid string) (*models.User, error) {
	var user models.User
	err := r.db.Where("uuid = ?", uuid).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", now).Error
}

func (r *UserRepository) FindVerificationByToken(token string) (*models.VerificationRequest, error) {
	var verification models.VerificationRequest
	err := r.db.Where("token = ?", token).First(&verification).Error
	return &verification, err
}

func (r *UserRepository) CreateVerification(verification *models.VerificationRequest) error {
	return r.db.Create(verification).Error
}

func (r *UserRepository) DeleteVerification(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.VerificationRequest{}).Error
}

func (r *UserRepository) DeleteVerificationByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.VerificationRequest{}).Error
}
