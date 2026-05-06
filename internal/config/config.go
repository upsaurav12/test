package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// jwtSecretMinLen is the minimum acceptable length for JWT_SECRET.
// A secret shorter than 32 bytes is brute-forceable offline against HS256.
const jwtSecretMinLen = 32

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

	// Database connection pool (tunable for horizontal scaling)
	DBMaxOpenConns        int
	DBMaxIdleConns        int
	DBConnMaxLifetimeSecs int

	// AutoMigrate enables GORM AutoMigrate on startup. Disable in production
	// and use versioned migration files via cmd/migrate instead.
	AutoMigrate bool

	// Auth
	JWTSecret      string
	JWTExpiryHours int

	// CORS – comma-separated list of allowed origins, e.g. "https://app.example.com"
	CORSAllowedOrigins []string
}

// New reads configuration from environment variables, applies defaults, and
// validates that all required fields are present.
func New() (*Config, error) {
	cfg := &Config{
		Port:                  getEnv("PORT", "8080"),
		AppEnv:                getEnv("APP_ENV", "production"),
		DBHost:                getEnv("GONE_DB_HOST", "localhost"),
		DBPort:                getEnv("GONE_DB_PORT", "5432"),
		DBName:                os.Getenv("GONE_DB_DATABASE"),
		DBUser:                os.Getenv("GONE_DB_USERNAME"),
		DBPassword:            os.Getenv("GONE_DB_PASSWORD"),
		DBSchema:              getEnv("GONE_DB_SCHEMA", "public"),
		DBSSLMode:             getEnv("GONE_DB_SSLMODE", "require"),
		DBMaxOpenConns:        getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:        getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetimeSecs: getEnvInt("DB_CONN_MAX_LIFETIME_SECS", 300),
		AutoMigrate:           getEnvBool("DB_AUTO_MIGRATE", false),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		JWTExpiryHours:        getEnvInt("JWT_EXPIRY_HOURS", 24),
		CORSAllowedOrigins:    strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
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
	var errs []string

	if c.DBName == "" {
		errs = append(errs, "GONE_DB_DATABASE is required")
	}
	if c.DBUser == "" {
		errs = append(errs, "GONE_DB_USERNAME is required")
	}
	if c.DBPassword == "" {
		errs = append(errs, "GONE_DB_PASSWORD is required")
	}
	if c.JWTSecret == "" {
		errs = append(errs, "JWT_SECRET is required")
	} else if len(c.JWTSecret) < jwtSecretMinLen {
		errs = append(errs, fmt.Sprintf("JWT_SECRET must be at least %d bytes (got %d)", jwtSecretMinLen, len(c.JWTSecret)))
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid configuration: %s", strings.Join(errs, "; "))
	}

	return nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return defaultVal
}
