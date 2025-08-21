package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/service"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body service.RegisterRequest true "User registration request"
// @Success 201 {object} service.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	authResponse, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, authResponse)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body service.LoginRequest true "User login credentials"
// @Success 200 {object} service.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	authResponse, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} entity.User
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	user, err := h.authService.GetUserProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile information of the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param profile body ProfileUpdateRequest true "Profile update request"
// @Success 200 {object} entity.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	var req ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Convert request to map
	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID.(uint), updates)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the password of the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param password body service.PasswordChangeRequest true "Password change request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	var req service.PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Password changed successfully",
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh an access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} service.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	tokenResponse, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// Logout godoc
// @Summary Logout user
// @Description Revoke user's access token
// @Tags auth
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid authorization header",
		})
		return
	}

	token := authHeader[7:] // Remove "Bearer " prefix

	err := h.authService.RevokeToken(c.Request.Context(), token)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Logged out successfully",
	})
}

// GetUser godoc
// @Summary Get user by ID (Admin only)
// @Description Get user information by ID (requires admin privileges)
// @Tags auth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entity.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/users/{id} [get]
func (h *AuthHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid user ID",
		})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// handleAuthError handles different types of authentication errors
func handleAuthError(c *gin.Context, err error) {
	switch err {
	case entity.ErrUserNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Not Found",
			Message: err.Error(),
		})
	case entity.ErrUserAlreadyExists:
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "Conflict",
			Message: err.Error(),
		})
	case entity.ErrInvalidCredentials:
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: err.Error(),
		})
	case entity.ErrUserInactive:
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: err.Error(),
		})
	case entity.ErrUnauthorized, entity.ErrInvalidToken:
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: err.Error(),
		})
	case entity.ErrUserEmailRequired, entity.ErrUserUsernameRequired, entity.ErrUserUsernameTooShort,
		 entity.ErrUserUsernameTooLong, entity.ErrInvalidInput, entity.ErrValidationFailed:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "An unexpected error occurred",
			Details: err.Error(),
		})
	}
}