package entity

import "errors"

// Product-related errors
var (
	ErrProductNotFound        = errors.New("product not found")
	ErrProductNameRequired    = errors.New("product name is required")
	ErrProductNameTooShort    = errors.New("product name must be at least 3 characters")
	ErrProductNameTooLong     = errors.New("product name must be less than 255 characters")
	ErrProductPriceInvalid    = errors.New("product price must be greater than or equal to 0")
	ErrProductStockInvalid    = errors.New("product stock must be greater than or equal to 0")
	ErrProductAlreadyExists   = errors.New("product with this name already exists")
)

// User-related errors
var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserEmailRequired      = errors.New("user email is required")
	ErrUserUsernameRequired   = errors.New("user username is required")
	ErrUserUsernameTooShort   = errors.New("username must be at least 3 characters")
	ErrUserUsernameTooLong    = errors.New("username must be less than 50 characters")
	ErrUserAlreadyExists      = errors.New("user with this email or username already exists")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrUserInactive           = errors.New("user account is inactive")
	ErrUnauthorized           = errors.New("unauthorized access")
	ErrInvalidToken           = errors.New("invalid or expired token")
)

// General errors
var (
	ErrInternalServer         = errors.New("internal server error")
	ErrInvalidInput           = errors.New("invalid input data")
	ErrDatabaseConnection     = errors.New("database connection error")
	ErrValidationFailed       = errors.New("validation failed")
)
