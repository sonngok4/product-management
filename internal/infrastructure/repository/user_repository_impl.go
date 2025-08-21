package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
	"gorm.io/gorm"
)

// userRepositoryImpl implements the UserRepository interface
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// Create creates a new user
func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by their ID
func (r *userRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email
func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetByUsername retrieves a user by their username
func (r *userRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// GetAll retrieves all users with optional filtering and pagination
func (r *userRepositoryImpl) GetAll(ctx context.Context, filter *repository.UserFilter, offset, limit int) ([]*entity.User, error) {
	var users []*entity.User
	query := r.db.WithContext(ctx)

	if filter != nil {
		query = r.applyFilter(query, filter)
	}

	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// GetTotalCount returns the total count of users with optional filtering
func (r *userRepositoryImpl) GetTotalCount(ctx context.Context, filter *repository.UserFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.User{})

	if filter != nil {
		query = r.applyFilter(query, filter)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// Update updates an existing user
func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete soft-deletes a user by their ID
func (r *userRepositoryImpl) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// HardDelete permanently deletes a user by their ID
func (r *userRepositoryImpl) HardDelete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(&entity.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to hard delete user: %w", err)
	}
	return nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence by email: %w", err)
	}
	return count > 0, nil
}

// ExistsByUsername checks if a user with the given username exists
func (r *userRepositoryImpl) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence by username: %w", err)
	}
	return count > 0, nil
}

// UpdateLastLogin updates the last login time for a user
func (r *userRepositoryImpl) UpdateLastLogin(ctx context.Context, id uint) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("last_login_at", &now).Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// UpdatePassword updates the password for a user
func (r *userRepositoryImpl) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// GetAdminUsers retrieves all admin users
func (r *userRepositoryImpl) GetAdminUsers(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.WithContext(ctx).Where("is_admin = ? AND is_active = ?", true, true).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin users: %w", err)
	}
	return users, nil
}

// applyFilter applies filters to the query
func (r *userRepositoryImpl) applyFilter(query *gorm.DB, filter *repository.UserFilter) *gorm.DB {
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	
	if filter.IsAdmin != nil {
		query = query.Where("is_admin = ?", *filter.IsAdmin)
	}
	
	if filter.SearchTerm != "" {
		searchPattern := "%" + filter.SearchTerm + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern)
	}
	
	return query
}
