package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port                string
	DatabaseURL         string
	RedisURL            string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	SendGridAPIKey      string
	ClerkPublishableKey string
	ClerkSecretKey      string
	StreamAPIKey        string
	StreamAPISecret     string
	AdminEmail          string
	AdminPassword       string
	AppEnv              string
}

// Load reads environment variables (from .env if present) and returns a populated Config.
// It panics if any required variable is missing.
func Load() *Config {
	// Load .env file if it exists; ignore error in production where vars may be injected directly.
	_ = godotenv.Load()

	cfg := &Config{
		Port:                os.Getenv("PORT"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		RedisURL:            os.Getenv("REDIS_URL"),
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
		SendGridAPIKey:      os.Getenv("SENDGRID_API_KEY"),
		ClerkPublishableKey: os.Getenv("CLERK_PUBLISHABLE_KEY"),
		ClerkSecretKey:      os.Getenv("CLERK_SECRET_KEY"),
		StreamAPIKey:        os.Getenv("STREAM_API_KEY"),
		StreamAPISecret:     os.Getenv("STREAM_API_SECRET"),
		AdminEmail:          os.Getenv("ADMIN_EMAIL"),
		AdminPassword:       os.Getenv("ADMIN_PASSWORD"),
		AppEnv:              os.Getenv("APP_ENV"),
	}

	// Validate required fields.
	if cfg.DatabaseURL == "" {
		panic(fmt.Errorf("config: DATABASE_URL is required but not set"))
	}
	if cfg.Port == "" {
		panic(fmt.Errorf("config: PORT is required but not set"))
	}

	return cfg
}
