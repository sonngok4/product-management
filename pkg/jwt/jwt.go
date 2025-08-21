package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/product-management/internal/domain/service"
)

// TokenManager handles JWT token operations
type TokenManager struct {
	secretKey string
	expiresIn time.Duration
}

// NewTokenManager creates a new token manager
func NewTokenManager(secretKey string, expiresIn time.Duration) *TokenManager {
	return &TokenManager{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

// CustomClaims represents the JWT claims
type CustomClaims struct {
	service.Claims
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token
func (tm *TokenManager) GenerateToken(claims *service.Claims) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(tm.expiresIn)

	customClaims := CustomClaims{
		Claims: *claims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   claims.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(tm.secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

// ValidateToken validates a JWT token and returns claims
func (tm *TokenManager) ValidateToken(tokenString string) (*service.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &claims.Claims, nil
}

// IsTokenExpired checks if a token is expired
func (tm *TokenManager) IsTokenExpired(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return true
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims.ExpiresAt.Before(time.Now())
	}

	return true
}

// GetTokenClaims extracts claims from a token without validation
func (tm *TokenManager) GetTokenClaims(tokenString string) (*service.Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &CustomClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return &claims.Claims, nil
	}

	return nil, errors.New("invalid token claims")
}