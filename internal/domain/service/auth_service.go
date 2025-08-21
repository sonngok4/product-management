package service

import (
	"context"

	"github.com/product-management/internal/domain/entity"
)

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User  *entity.User   `json:"user"`
	Token *TokenResponse `json:"token"`
}

// PasswordChangeRequest represents a password change request
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

// AuthService defines the interface for authentication business logic operations
type AuthService interface {
	// Register creates a new user account
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	
	// Login authenticates a user and returns tokens
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	
	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	
	// UpdateProfile updates user profile information
	UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) (*entity.User, error)
	
	// ChangePassword changes user password
	ChangePassword(ctx context.Context, userID uint, req *PasswordChangeRequest) error
	
	// GenerateToken generates a new JWT token for the user
	GenerateToken(ctx context.Context, user *entity.User) (*TokenResponse, error)
	
	// ValidateToken validates a JWT token and returns claims
	ValidateToken(ctx context.Context, token string) (*Claims, error)
	
	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	
	// RevokeToken revokes a token (logout)
	RevokeToken(ctx context.Context, token string) error
	
	// GetUserProfile gets user profile information
	GetUserProfile(ctx context.Context, userID uint) (*entity.User, error)
}
