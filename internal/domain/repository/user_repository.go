package repository

import (
	"context"

	"github.com/product-management/internal/domain/entity"
)

// UserFilter represents filtering criteria for users
type UserFilter struct {
	IsActive    *bool
	IsAdmin     *bool
	SearchTerm  string // for searching in username, email, first_name, or last_name
}

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error
	
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	
	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	
	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	
	// GetAll retrieves all users with optional filtering and pagination
	GetAll(ctx context.Context, filter *UserFilter, offset, limit int) ([]*entity.User, error)
	
	// GetTotalCount returns the total count of users with optional filtering
	GetTotalCount(ctx context.Context, filter *UserFilter) (int64, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error
	
	// Delete soft-deletes a user by their ID
	Delete(ctx context.Context, id uint) error
	
	// HardDelete permanently deletes a user by their ID
	HardDelete(ctx context.Context, id uint) error
	
	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	
	// ExistsByUsername checks if a user with the given username exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	
	// UpdateLastLogin updates the last login time for a user
	UpdateLastLogin(ctx context.Context, id uint) error
	
	// UpdatePassword updates the password for a user
	UpdatePassword(ctx context.Context, id uint, hashedPassword string) error
	
	// GetAdminUsers retrieves all admin users
	GetAdminUsers(ctx context.Context) ([]*entity.User, error)
}