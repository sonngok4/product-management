package unit

import (
	"testing"
	"time"

	"github.com/product-management/internal/domain/service"
	"github.com/product-management/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// JWTTestSuite represents the test suite for JWT token manager
type JWTTestSuite struct {
	suite.Suite
	tokenManager *jwt.TokenManager
	secretKey    string
	expiresIn    time.Duration
}

// SetupTest sets up the test suite
func (suite *JWTTestSuite) SetupTest() {
	suite.secretKey = "test-secret-key"
	suite.expiresIn = time.Hour
	suite.tokenManager = jwt.NewTokenManager(suite.secretKey, suite.expiresIn)
}

// TestGenerateToken_Success tests successful token generation
func (suite *JWTTestSuite) TestGenerateToken_Success() {
	// Arrange
	claims := &service.Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	// Act
	token, expiresAt, err := suite.tokenManager.GenerateToken(claims)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	assert.True(suite.T(), expiresAt > time.Now().Unix())
}

// TestValidateToken_Success tests successful token validation
func (suite *JWTTestSuite) TestValidateToken_Success() {
	// Arrange
	originalClaims := &service.Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  true,
	}

	token, _, err := suite.tokenManager.GenerateToken(originalClaims)
	assert.NoError(suite.T(), err)

	// Act
	validatedClaims, err := suite.tokenManager.ValidateToken(token)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), validatedClaims)
	assert.Equal(suite.T(), originalClaims.UserID, validatedClaims.UserID)
	assert.Equal(suite.T(), originalClaims.Username, validatedClaims.Username)
	assert.Equal(suite.T(), originalClaims.Email, validatedClaims.Email)
	assert.Equal(suite.T(), originalClaims.IsAdmin, validatedClaims.IsAdmin)
}

// TestValidateToken_InvalidToken tests token validation with invalid token
func (suite *JWTTestSuite) TestValidateToken_InvalidToken() {
	// Arrange
	invalidToken := "invalid.token.here"

	// Act
	claims, err := suite.tokenManager.ValidateToken(invalidToken)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

// TestValidateToken_ExpiredToken tests token validation with expired token
func (suite *JWTTestSuite) TestValidateToken_ExpiredToken() {
	// Arrange - create a token manager with very short expiry
	shortExpiryTokenManager := jwt.NewTokenManager(suite.secretKey, time.Millisecond)
	
	claims := &service.Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	token, _, err := shortExpiryTokenManager.GenerateToken(claims)
	assert.NoError(suite.T(), err)

	// Wait for token to expire
	time.Sleep(time.Millisecond * 10)

	// Act
	validatedClaims, err := suite.tokenManager.ValidateToken(token)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
}

// TestValidateToken_WrongSecret tests token validation with wrong secret
func (suite *JWTTestSuite) TestValidateToken_WrongSecret() {
	// Arrange
	claims := &service.Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	token, _, err := suite.tokenManager.GenerateToken(claims)
	assert.NoError(suite.T(), err)

	// Create token manager with different secret
	wrongSecretTokenManager := jwt.NewTokenManager("wrong-secret", suite.expiresIn)

	// Act
	validatedClaims, err := wrongSecretTokenManager.ValidateToken(token)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), validatedClaims)
}

// TestIsTokenExpired_ValidToken tests checking expiry of valid token
func (suite *JWTTestSuite) TestIsTokenExpired_ValidToken() {
	// Arrange
	claims := &service.Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	token, _, err := suite.tokenManager.GenerateToken(claims)
	assert.NoError(suite.T(), err)

	// Act
	isExpired := suite.tokenManager.IsTokenExpired(token)

	// Assert
	assert.False(suite.T(), isExpired)
}

// TestIsTokenExpired_ExpiredToken tests checking expiry of expired token
func (suite *JWTTestSuite) TestIsTokenExpired_ExpiredToken() {
	// Arrange - create a token manager with very short expiry
	shortExpiryTokenManager := jwt.NewTokenManager(suite.secretKey, time.Millisecond)
	
	claims := &service.Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	token, _, err := shortExpiryTokenManager.GenerateToken(claims)
	assert.NoError(suite.T(), err)

	// Wait for token to expire
	time.Sleep(time.Millisecond * 10)

	// Act
	isExpired := suite.tokenManager.IsTokenExpired(token)

	// Assert
	assert.True(suite.T(), isExpired)
}

// TestIsTokenExpired_InvalidToken tests checking expiry of invalid token
func (suite *JWTTestSuite) TestIsTokenExpired_InvalidToken() {
	// Arrange
	invalidToken := "invalid.token.here"

	// Act
	isExpired := suite.tokenManager.IsTokenExpired(invalidToken)

	// Assert
	assert.True(suite.T(), isExpired) // Invalid tokens are considered expired
}

// TestGetTokenClaims_Success tests extracting claims without validation
func (suite *JWTTestSuite) TestGetTokenClaims_Success() {
	// Arrange
	originalClaims := &service.Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  true,
	}

	token, _, err := suite.tokenManager.GenerateToken(originalClaims)
	assert.NoError(suite.T(), err)

	// Act
	extractedClaims, err := suite.tokenManager.GetTokenClaims(token)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), extractedClaims)
	assert.Equal(suite.T(), originalClaims.UserID, extractedClaims.UserID)
	assert.Equal(suite.T(), originalClaims.Username, extractedClaims.Username)
	assert.Equal(suite.T(), originalClaims.Email, extractedClaims.Email)
	assert.Equal(suite.T(), originalClaims.IsAdmin, extractedClaims.IsAdmin)
}

// TestGetTokenClaims_InvalidToken tests extracting claims from invalid token
func (suite *JWTTestSuite) TestGetTokenClaims_InvalidToken() {
	// Arrange
	invalidToken := "invalid.token.here"

	// Act
	claims, err := suite.tokenManager.GetTokenClaims(invalidToken)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

// TestJWTTestSuite runs the test suite
func TestJWTTestSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}