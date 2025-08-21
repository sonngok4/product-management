package handler

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// StockUpdateRequest represents a request to update product stock
type StockUpdateRequest struct {
	Stock int `json:"stock" validate:"min=0"`
}

// ProfileUpdateRequest represents a request to update user profile
type ProfileUpdateRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// BulkUpdateStatusRequest represents a request to bulk update product status
type BulkUpdateStatusRequest struct {
	ProductIDs []uint `json:"product_ids" validate:"required,min=1"`
	IsActive   bool   `json:"is_active"`
}
