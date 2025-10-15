package repositories

import (
	"time"

	"github.com/iqmalr-pedia/go-auth/v2/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}
type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(id uint) (*models.User, error)
	FindUserByUUID(uuid string) (*models.User, error)
	FindUserByPhone(phone string) (*models.User, error)
	UpdateUser(user *models.User) error
	UpdateLastLogin(userID uint) error
	CreatePasswordResetToken(token *models.PasswordResetToken) error
	FindPasswordResetToken(token string) (*models.PasswordResetToken, error)
	UpdatePasswordResetToken(token *models.PasswordResetToken) error
	DeleteExpiredPasswordResetTokens() error
	CreateEmailVerificationToken(token *models.EmailVerificationToken) error
	FindEmailVerificationToken(token string) (*models.EmailVerificationToken, error)
	DeleteEmailVerificationToken(token string) error
	DeleteExpiredEmailVerificationTokens() error
	FindEmailVerificationTokenByUserID(userID uint) (*models.EmailVerificationToken, error)
	DeleteEmailVerificationTokenByUserID(userID uint) error
	UpdateProfile(userID uint, updates map[string]interface{}) error
	ChangePassword(userID uint, hashedPassword string) error
	FindUsersWithPagination(page, limit int, search, role string, isActive *bool) ([]models.User, int64, error)
	DeactivateUser(userID uint) error
	FindUserByEmailAndNotID(email string, id uint) (*models.User, error)
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

func (r *UserRepository) CreatePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *UserRepository) FindPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&resetToken).Error
	return &resetToken, err
}

func (r *UserRepository) UpdatePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Save(token).Error
}

func (r *UserRepository) DeleteExpiredPasswordResetTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{}).Error
}

func (r *UserRepository) CreateEmailVerificationToken(token *models.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

func (r *UserRepository) FindEmailVerificationToken(token string) (*models.EmailVerificationToken, error) {
	var verificationToken models.EmailVerificationToken
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&verificationToken).Error
	return &verificationToken, err
}

func (r *UserRepository) DeleteEmailVerificationToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.EmailVerificationToken{}).Error
}

func (r *UserRepository) DeleteExpiredEmailVerificationTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.EmailVerificationToken{}).Error
}

func (r *UserRepository) FindEmailVerificationTokenByUserID(userID uint) (*models.EmailVerificationToken, error) {
	var verificationToken models.EmailVerificationToken
	err := r.db.Where("user_id = ?", userID).First(&verificationToken).Error
	return &verificationToken, err
}

func (r *UserRepository) DeleteEmailVerificationTokenByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.EmailVerificationToken{}).Error
}

func (r *UserRepository) UpdateProfile(userID uint, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *UserRepository) ChangePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

func (r *UserRepository) FindUsersWithPagination(page, limit int, search, role string, isActive *bool) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) DeactivateUser(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error
}

func (r *UserRepository) FindUserByEmailAndNotID(email string, id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND id != ?", email, id).First(&user).Error
	return &user, err
}
