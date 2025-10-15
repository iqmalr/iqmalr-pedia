package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iqmalr-pedia/go-auth/v2/internal/dto/request"
	"github.com/iqmalr-pedia/go-auth/v2/internal/dto/response"
	"github.com/iqmalr-pedia/go-auth/v2/internal/models"
	"github.com/iqmalr-pedia/go-auth/v2/internal/repositories"
	"github.com/iqmalr-pedia/go-auth/v2/internal/utils"
)

type AuthService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewAuthService(userRepo repositories.UserRepositoryInterface) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req *request.RegisterRequest) (*response.RegisterResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.FindUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		UUID:     uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
	}

	if user.Role == "" {
		user.Role = "customer"
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Create email verification token
	token := &models.EmailVerificationToken{
		UserID:    user.ID,
		Token:     utils.GenerateRandomToken(32),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.userRepo.CreateEmailVerificationToken(token)
	if err != nil {
		return nil, err
	}

	// TODO: Send verification email

	return &response.RegisterResponse{
		Message: "User registered successfully. Please check your email for verification.",
		User: &response.UserResponse{
			ID:        user.ID,
			UUID:      user.UUID.String(),
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Update last login
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.UUID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.UUID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Save refresh token to user
	user.RefreshToken = refreshToken
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Message: "Login successful",
		Token:   token,
		User: &response.LoginUser{
			ID:       user.ID,
			UUID:     user.UUID.String(),
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, nil
}

func (s *AuthService) Logout(userID uint) error {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return err
	}

	// Clear refresh token
	user.RefreshToken = ""
	return s.userRepo.UpdateUser(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*response.RefreshTokenResponse, error) {
	// Verify refresh token
	claims, err := utils.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Find user
	user, err := s.userRepo.FindUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if refresh token matches
	if user.RefreshToken != refreshToken {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new access token
	token, err := utils.GenerateToken(user.ID, user.UUID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &response.RefreshTokenResponse{
		Token: token,
	}, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return errors.New("user with this email does not exist")
	}

	// Create password reset token
	token := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     utils.GenerateRandomToken(32),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	err = s.userRepo.CreatePasswordResetToken(token)
	if err != nil {
		return err
	}

	// TODO: Send password reset email

	return nil
}

func (s *AuthService) ResetPassword(req *request.ResetPasswordRequest) error {
	// Find reset token
	resetToken, err := s.userRepo.FindPasswordResetToken(req.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// Find user
	user, err := s.userRepo.FindUserByID(resetToken.UserID)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	// Mark token as used
	resetToken.Used = true
	err = s.userRepo.UpdatePasswordResetToken(resetToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyEmail(token string) error {
	// Find verification token
	verificationToken, err := s.userRepo.FindEmailVerificationToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// Find user
	user, err := s.userRepo.FindUserByID(verificationToken.UserID)
	if err != nil {
		return err
	}

	// Update email verification
	now := time.Now()
	user.EmailVerifiedAt = &now
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	// Delete verification token
	err = s.userRepo.DeleteEmailVerificationToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResendVerification(email string) error {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return errors.New("user with this email does not exist")
	}

	if user.EmailVerifiedAt != nil {
		return errors.New("email is already verified")
	}

	// Delete existing verification tokens
	err = s.userRepo.DeleteEmailVerificationTokenByUserID(user.ID)
	if err != nil {
		return err
	}

	// Create new verification token
	token := &models.EmailVerificationToken{
		UserID:    user.ID,
		Token:     utils.GenerateRandomToken(32),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.userRepo.CreateEmailVerificationToken(token)
	if err != nil {
		return err
	}

	// TODO: Send verification email

	return nil
}

func (s *AuthService) GetUserProfile(userID uint) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &response.UserProfileResponse{
		ID:              user.ID,
		UUID:            user.UUID.String(),
		Name:            user.Name,
		Email:           user.Email,
		Phone:           user.Phone,
		Role:            user.Role,
		IsActive:        user.IsActive,
		EmailVerifiedAt: user.EmailVerifiedAt,
		PhoneVerifiedAt: user.PhoneVerifiedAt,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
	}, nil
}
