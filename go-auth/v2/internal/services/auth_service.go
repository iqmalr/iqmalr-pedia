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

type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *AuthService) Register(req *request.RegisterRequest) (*response.RegisterResponse, error) {
	_, err := s.userRepo.FindUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

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
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.UUID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.UUID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

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

	user.RefreshToken = ""
	return s.userRepo.UpdateUser(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*response.RefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.RefreshToken != refreshToken {
		return nil, errors.New("invalid refresh token")
	}

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
	resetToken, err := s.userRepo.FindPasswordResetToken(req.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	user, err := s.userRepo.FindUserByID(resetToken.UserID)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	resetToken.Used = true
	err = s.userRepo.UpdatePasswordResetToken(resetToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyEmail(token string) error {
	verificationToken, err := s.userRepo.FindEmailVerificationToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	user, err := s.userRepo.FindUserByID(verificationToken.UserID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

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

	err = s.userRepo.DeleteEmailVerificationTokenByUserID(user.ID)
	if err != nil {
		return err
	}

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
func (s *UserService) GetProfile(userID uint) (*response.UserProfileResponse, error) {
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
		AvatarUrl:       user.AvatarUrl,
		IsActive:        user.IsActive,
		EmailVerifiedAt: user.EmailVerifiedAt,
		PhoneVerifiedAt: user.PhoneVerifiedAt,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}, nil
}

func (s *UserService) UpdateProfile(userID uint, req *request.UpdateProfileRequest) (*response.UserProfileResponse, error) {
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.AvatarUrl != "" {
		updates["avatar_url"] = req.AvatarUrl
	}

	if len(updates) == 0 {
		return s.GetProfile(userID)
	}

	if err := s.userRepo.UpdateProfile(userID, updates); err != nil {
		return nil, err
	}

	return s.GetProfile(userID)
}

func (s *UserService) ChangePassword(userID uint, req *request.ChangePasswordRequest) error {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.ChangePassword(userID, hashedPassword)
}

func (s *UserService) GetUserByID(id uint) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.FindUserByID(id)
	if err != nil {
		return nil, err
	}
	return &response.UserProfileResponse{
		ID:              user.ID,
		UUID:            user.UUID.String(),
		Email:           user.Email,
		Phone:           user.Phone,
		Role:            user.Role,
		AvatarUrl:       user.AvatarUrl,
		IsActive:        user.IsActive,
		EmailVerifiedAt: user.EmailVerifiedAt,
		PhoneVerifiedAt: user.PhoneVerifiedAt,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}, nil
}

func (s *UserService) UpdateUser(id uint, req *request.UpdateUserRequest) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.FindUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Email != "" && req.Email != user.Email {
		_, err := s.userRepo.FindUserByEmailAndNotID(req.Email, id)
		if err == nil {
			return nil, errors.New("email is already taken by another user")
		}
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.userRepo.UpdateProfile(id, updates); err != nil {
		return nil, err
	}

	return s.GetUserByID(id)
}

func (s *UserService) DeactivateUser(id uint) error {
	if err := s.userRepo.DeactivateUser(id); err != nil {
		return errors.New("failed to deactivate user")
	}
	return nil
}

func (s *UserService) ListUsers(req *request.ListUsersRequest) (*response.UserListResponse, error) {
	users, total, err := s.userRepo.FindUsersWithPagination(req.Page, req.Limit, req.Search, req.Role, req.IsActive)
	if err != nil {
		return nil, err
	}

	userListItems := make([]response.UserListItem, len(users))
	for i, user := range users {
		userListItems[i] = response.UserListItem{
			ID:        user.ID,
			UUID:      user.UUID.String(),
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		}
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &response.UserListResponse{
		Data: userListItems,
		Pagination: response.PaginationResponse{
			TotalPages: totalPages,
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
		},
	}, nil
}
