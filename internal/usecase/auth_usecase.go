package usecase

import (
	"context"
	"regexp"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
	"github.com/product-management/internal/domain/service"
	"github.com/product-management/pkg/jwt"
)

// authUseCase implements the AuthService interface
type authUseCase struct {
	userRepo     repository.UserRepository
	tokenManager *jwt.TokenManager
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(userRepo repository.UserRepository, tokenManager *jwt.TokenManager) service.AuthService {
	return &authUseCase{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

// Register creates a new user account
func (uc *authUseCase) Register(ctx context.Context, req *service.RegisterRequest) (*service.AuthResponse, error) {
	// Validate input
	if err := uc.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if user already exists by email
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entity.ErrUserAlreadyExists
	}

	// Check if user already exists by username
	exists, err = uc.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entity.ErrUserAlreadyExists
	}

	// Create user entity
	user := &entity.User{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		IsAdmin:   false,
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := uc.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log but don't fail the registration
	}

	return &service.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// Login authenticates a user and returns tokens
func (uc *authUseCase) Login(ctx context.Context, req *service.LoginRequest) (*service.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, entity.ErrInvalidCredentials
	}

	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == entity.ErrUserNotFound {
			return nil, entity.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, entity.ErrUserInactive
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, entity.ErrInvalidCredentials
	}

	// Generate token
	token, err := uc.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log but don't fail the login
	}

	return &service.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by their ID
func (uc *authUseCase) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

// UpdateProfile updates user profile information
func (uc *authUseCase) UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) (*entity.User, error) {
	// Get current user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if firstName, ok := updates["first_name"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := updates["last_name"].(string); ok {
		user.LastName = lastName
	}
	if username, ok := updates["username"].(string); ok {
		if username != user.Username {
			// Check if username is already taken
			exists, err := uc.userRepo.ExistsByUsername(ctx, username)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, entity.ErrUserAlreadyExists
			}
			user.Username = username
		}
	}

	// Validate updated user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Save updates
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes user password
func (uc *authUseCase) ChangePassword(ctx context.Context, userID uint, req *service.PasswordChangeRequest) error {
	// Get current user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := user.CheckPassword(req.CurrentPassword); err != nil {
		return entity.ErrInvalidCredentials
	}

	// Validate new password
	if err := uc.validatePassword(req.NewPassword); err != nil {
		return err
	}

	// Hash new password
	if err := user.HashPassword(req.NewPassword); err != nil {
		return err
	}

	// Update password in database
	return uc.userRepo.UpdatePassword(ctx, userID, user.Password)
}

// GenerateToken generates a new JWT token for the user
func (uc *authUseCase) GenerateToken(ctx context.Context, user *entity.User) (*service.TokenResponse, error) {
	claims := &service.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}

	tokenString, expiresAt, err := uc.tokenManager.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	return &service.TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expiresAt,
	}, nil
}

// ValidateToken validates a JWT token and returns claims
func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*service.Claims, error) {
	return uc.tokenManager.ValidateToken(token)
}

// RefreshToken refreshes an access token using a refresh token
func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*service.TokenResponse, error) {
	// For now, we'll treat the refresh token as a regular token
	// In a production system, you'd have separate refresh token logic
	claims, err := uc.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, entity.ErrInvalidToken
	}

	// Get user to ensure they still exist and are active
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, entity.ErrUserInactive
	}

	// Generate new token
	return uc.GenerateToken(ctx, user)
}

// RevokeToken revokes a token (logout)
func (uc *authUseCase) RevokeToken(ctx context.Context, token string) error {
	// In a production system, you'd maintain a blacklist of revoked tokens
	// For now, we'll just validate the token to ensure it's valid
	_, err := uc.tokenManager.ValidateToken(token)
	return err
}

// GetUserProfile gets user profile information
func (uc *authUseCase) GetUserProfile(ctx context.Context, userID uint) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

// validateRegisterRequest validates registration request
func (uc *authUseCase) validateRegisterRequest(req *service.RegisterRequest) error {
	if req.Email == "" {
		return entity.ErrUserEmailRequired
	}

	if !uc.isValidEmail(req.Email) {
		return entity.ErrInvalidInput
	}

	if req.Username == "" {
		return entity.ErrUserUsernameRequired
	}

	if len(req.Username) < 3 {
		return entity.ErrUserUsernameTooShort
	}

	if len(req.Username) > 50 {
		return entity.ErrUserUsernameTooLong
	}

	return uc.validatePassword(req.Password)
}

// validatePassword validates password strength
func (uc *authUseCase) validatePassword(password string) error {
	if len(password) < 8 {
		return entity.ErrInvalidInput
	}
	return nil
}

// isValidEmail validates email format
func (uc *authUseCase) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}