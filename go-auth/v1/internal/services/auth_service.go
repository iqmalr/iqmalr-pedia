package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iqmalr-pedia/go-auth/v1/internal/dto/request"
	"github.com/iqmalr-pedia/go-auth/v1/internal/dto/response"
	"github.com/iqmalr-pedia/go-auth/v1/internal/models"
	"github.com/iqmalr-pedia/go-auth/v1/internal/repositories"
	"github.com/iqmalr-pedia/go-auth/v1/internal/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req *request.RegisterRequest) (*response.AuthResponse, error) {
	existingUser, _ := s.userRepo.FindUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	if req.Phone != "" {
		existingPhoneUser, _ := s.userRepo.FindUserByPhone(req.Phone)
		if existingPhoneUser != nil {
			return nil, errors.New("user with this phone number already exists")
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     "customer",
		IsActive: true,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	verificationToken := uuid.New().String()
	verification := &models.VerificationRequest{
		UserID:    user.ID,
		Token:     verificationToken,
		Type:      "email",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.userRepo.CreateVerification(verification); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.UUID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	user.RefreshToken = refreshToken
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		UserID:        user.ID,
		UUID:          user.UUID.String(),
		Name:          user.Name,
		Email:         user.Email,
		Phone:         user.Phone,
		Role:          user.Role,
		AvatarURL:     user.AvatarURL,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerifiedAt != nil,
		PhoneVerified: user.PhoneVerifiedAt != nil,
	}, nil
}
func (s *AuthService) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.UUID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	user.RefreshToken = refreshToken
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		UserID:        user.ID,
		UUID:          user.UUID.String(),
		Name:          user.Name,
		Email:         user.Email,
		Phone:         user.Phone,
		Role:          user.Role,
		AvatarURL:     user.AvatarURL,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerifiedAt != nil,
		PhoneVerified: user.PhoneVerifiedAt != nil,
	}, nil
}

func (s *AuthService) Verify(token string, verifyType string) error {
	verification, err := s.userRepo.FindVerificationByToken(token)
	if err != nil {
		return errors.New("invalid verification token")
	}

	if verification.Type != verifyType {
		return errors.New("invalid verification type")
	}

	if time.Now().After(verification.ExpiresAt) {
		return errors.New("verification token expired")
	}

	user, err := s.userRepo.FindUserByID(verification.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	now := time.Now()
	if verifyType == "email" {
		user.EmailVerifiedAt = &now
	} else if verifyType == "phone" {
		user.PhoneVerifiedAt = &now
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	err = s.userRepo.DeleteVerification(token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshToken(req *request.RefreshRequest) (*response.AuthResponse, error) {
	claims, err := utils.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindUserByID(claims.UserID)
	if err != nil || user.RefreshToken != req.RefreshToken {
		return nil, errors.New("invalid refresh token")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.UUID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	user.RefreshToken = refreshToken
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		UserID:        user.ID,
		UUID:          user.UUID.String(),
		Name:          user.Name,
		Email:         user.Email,
		Phone:         user.Phone,
		Role:          user.Role,
		AvatarURL:     user.AvatarURL,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerifiedAt != nil,
		PhoneVerified: user.PhoneVerifiedAt != nil,
	}, nil
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
		AvatarURL:       user.AvatarURL,
		IsActive:        user.IsActive,
		EmailVerifiedAt: user.EmailVerifiedAt,
		PhoneVerifiedAt: user.PhoneVerifiedAt,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
	}, nil
}

func (s *AuthService) UpdateUserProfile(userID uint, req *request.UpdateProfileRequest) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Phone != "" && req.Phone != user.Phone {
		existingPhoneUser, _ := s.userRepo.FindUserByPhone(req.Phone)
		if existingPhoneUser != nil {
			return nil, errors.New("phone number already in use")
		}
		user.Phone = req.Phone
		user.PhoneVerifiedAt = nil
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return s.GetUserProfile(userID)
}

func (s *AuthService) Logout(userID uint) error {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return err
	}

	user.RefreshToken = ""
	return s.userRepo.UpdateUser(user)
}
