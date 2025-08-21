package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth2   OAuth2Config
	CORS     CORSConfig
	Log      LogConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port    string
	GinMode string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret    string
	ExpiresIn string
}

// OAuth2Config holds OAuth2 configuration
type OAuth2Config struct {
	Google GoogleOAuth2Config
}

// GoogleOAuth2Config holds Google OAuth2 configuration
type GoogleOAuth2Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "product_management"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your-secret-key"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		},
		OAuth2: OAuth2Config{
			Google: GoogleOAuth2Config{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
			},
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvAsSlice("ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	return config
}

// Helper functions for environment variable handling
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(name string, defaultValue []string) []string {
	valueStr := getEnv(name, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}