package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user entity in the domain layer
type User struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Email       string         `json:"email" gorm:"uniqueIndex;size:255;not null" validate:"required,email"`
	Username    string         `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50"`
	Password    string         `json:"-" gorm:"size:255;not null"`
	FirstName   string         `json:"first_name" gorm:"size:100"`
	LastName    string         `json:"last_name" gorm:"size:100"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	IsAdmin     bool           `json:"is_admin" gorm:"default:false"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for User entity
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// HashPassword hashes the user's password
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the user's password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}

// Validate performs basic validation on the user entity
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrUserEmailRequired
	}
	if u.Username == "" {
		return ErrUserUsernameRequired
	}
	if len(u.Username) < 3 {
		return ErrUserUsernameTooShort
	}
	if len(u.Username) > 50 {
		return ErrUserUsernameTooLong
	}
	return nil
}