package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Server
	Port   string
	AppEnv string

	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSchema   string
	DBSSLMode  string

	// Auth
	JWTSecret     string
	JWTExpiryHours int

	// CORS – comma-separated list of allowed origins, e.g. "https://app.example.com"
	CORSAllowedOrigins []string
}

// New reads configuration from environment variables, applies defaults, and
// validates that all required fields are present.
func New() (*Config, error) {
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "production"),
		DBHost:             getEnv("GONE_DB_HOST", "localhost"),
		DBPort:             getEnv("GONE_DB_PORT", "5432"),
		DBName:             os.Getenv("GONE_DB_DATABASE"),
		DBUser:             os.Getenv("GONE_DB_USERNAME"),
		DBPassword:         os.Getenv("GONE_DB_PASSWORD"),
		DBSchema:           getEnv("GONE_DB_SCHEMA", "public"),
		DBSSLMode:          getEnv("GONE_DB_SSLMODE", "require"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpiryHours:     24,
		CORSAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsDevelopment returns true when running in a local/development environment.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "local" || c.AppEnv == "development"
}

func (c *Config) validate() error {
	var errs []error

	if c.DBName == "" {
		errs = append(errs, errors.New("GONE_DB_DATABASE is required"))
	}
	if c.DBUser == "" {
		errs = append(errs, errors.New("GONE_DB_USERNAME is required"))
	}
	if c.DBPassword == "" {
		errs = append(errs, errors.New("GONE_DB_PASSWORD is required"))
	}
	if c.JWTSecret == "" {
		errs = append(errs, errors.New("JWT_SECRET is required"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("configuration errors: %w", errors.Join(errs...))
	}

	return nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
