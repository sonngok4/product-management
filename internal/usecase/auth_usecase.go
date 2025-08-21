package usecase

import (
	"context"
	"errors"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
	"github.com/product-management/pkg/jwt"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo     repository.UserRepository
	tokenManager *jwt.TokenManager
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(userRepo repository.UserRepository, tokenManager *jwt.TokenManager) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

// LoginRequest represents login request data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response data
type LoginResponse struct {
	Token string      `json:"token"`
	User  *entity.User `json:"user"`
}

// Login authenticates a user and returns a JWT token
func (uc *AuthUseCase) Login(req *LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := uc.userRepo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Determine role based on IsAdmin field
	role := "user"
	if user.IsAdmin {
		role = "admin"
	}

	// Generate JWT token
	token, err := uc.tokenManager.GenerateToken(user.ID, user.Email, role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// RegisterRequest represents registration request data
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
}

// Register creates a new user account
func (uc *AuthUseCase) Register(req *RegisterRequest) (*entity.User, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(context.Background(), req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Create user
	user := &entity.User{
		Email:     req.Email,
		Username:  req.Name, // Use Name as Username
		FirstName: req.Name, // Use Name as FirstName
		IsActive:  true,
		IsAdmin:   req.Role == "admin",
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}
