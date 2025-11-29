package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Gemini       GeminiConfig
	Zoom         ZoomConfig
	Drive        DriveConfig
	Notification NotificationConfig
	User         UserConfig
	CDC          CDCConfig
	CORS         CORSConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string
}

// GeminiConfig holds Gemini API configuration
type GeminiConfig struct {
	APIKey string
}

// ZoomConfig holds Zoom API configuration
type ZoomConfig struct {
	APIKey    string
	APISecret string
}

// DriveConfig holds Google Drive API configuration
type DriveConfig struct {
	APIKey string
}

// NotificationConfig holds notification service configuration
type NotificationConfig struct {
	APIKey string
}

// UserConfig holds user service API configuration
type UserConfig struct {
	BaseURL string
	APIKey  string
}

// CDCConfig holds CDC API configuration
type CDCConfig struct {
	BaseURL    string
	WebBaseURL string
	APIKey     string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Build DSN from individual environment variables
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "")
	dbName := getEnv("POSTGRES_DB", "moscow")
	sslMode := getEnv("POSTGRES_SSLMODE", "disable")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslMode)

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "5002"),
		},
		Database: DatabaseConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			DBName:   dbName,
			SSLMode:  sslMode,
			DSN:      dsn,
		},
		Gemini: GeminiConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		},
		Zoom: ZoomConfig{
			APIKey:    os.Getenv("ZOOM_API_KEY"),
			APISecret: os.Getenv("ZOOM_API_SECRET"),
		},
		Drive: DriveConfig{
			APIKey: os.Getenv("GOOGLE_DRIVE_API_KEY"),
		},
		Notification: NotificationConfig{
			APIKey: os.Getenv("NOTIFICATION_API_KEY"),
		},
		User: UserConfig{
			BaseURL: getEnv("USER_SERVICE_BASE_URL", "http://localhost:5001/api/v1/external"),
			APIKey:  getEnv("USER_SERVICE_API_KEY", "56c290ad131b1f3e3131059c6c33ff46be0cff5cab3673de2bf2c1d81798b1d8"),
		},
		CDC: CDCConfig{
			BaseURL:    getEnv("CDC_API_BASE_URL", "https://travel.state.gov/_travel-resources/content/travel-resources/www.tripsofia.com/api/v1"),
			WebBaseURL: getEnv("CDC_WEB_BASE_URL", "https://wwwnc.cdc.gov"),
			APIKey:     getEnv("CDC_API_KEY", ""),
		},
		CORS: CORSConfig{
			AllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Database configuration validation
	if c.Database.User == "" {
		log.Println("⚠️  WARNING: POSTGRES_USER not set")
	}
	if c.Database.Password == "" {
		log.Println("⚠️  WARNING: POSTGRES_PASSWORD not set")
	}
	if c.Database.DBName == "" {
		log.Println("⚠️  WARNING: POSTGRES_DB not set")
	}

	// Log database connection info (without password)
	log.Printf("📊 Database Config: Host=%s, Port=%s, User=%s, DB=%s, SSL=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.DBName, c.Database.SSLMode)

	// Gemini API Key is optional for basic functionality
	// If not provided, transaction extraction won't work but other features will
	if c.Gemini.APIKey == "" {
		log.Println("⚠️  WARNING: GEMINI_API_KEY not set - transaction extraction will not work")
	}

	// Optional validation for meeting functionality
	if c.Zoom.APIKey == "" {
		// Log warning but don't fail - Zoom functionality won't work
		// In production, you might want to make this required
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
