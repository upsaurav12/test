package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validSecret is exactly 32 bytes – the minimum accepted length.
const validSecret = "abcdefghijklmnopqrstuvwxyz123456"

// setRequiredEnv sets the minimum env vars needed for a valid Config.
func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("GONE_DB_DATABASE", "testdb")
	t.Setenv("GONE_DB_USERNAME", "testuser")
	t.Setenv("GONE_DB_PASSWORD", "testpass")
	t.Setenv("JWT_SECRET", validSecret)
}

func TestNew_PortFromEnv(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("PORT", "9090")

	cfg, err := New()

	require.NoError(t, err)
	assert.Equal(t, "9090", cfg.Port)
}

func TestNew_DefaultPortWhenEnvMissing(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("PORT", "")

	cfg, err := New()

	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.Port)
}

func TestNew_DefaultSSLMode(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := New()

	require.NoError(t, err)
	assert.Equal(t, "require", cfg.DBSSLMode)
}

func TestNew_MissingRequiredFields(t *testing.T) {
	t.Setenv("GONE_DB_DATABASE", "")
	t.Setenv("GONE_DB_USERNAME", "")
	t.Setenv("GONE_DB_PASSWORD", "")
	t.Setenv("JWT_SECRET", "")

	_, err := New()

	assert.Error(t, err)
}

func TestNew_IsDevelopment(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("APP_ENV", "local")

	cfg, err := New()

	require.NoError(t, err)
	assert.True(t, cfg.IsDevelopment())
}

func TestNew_JWTSecretTooShort(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("JWT_SECRET", "tooshort") // only 8 bytes

	_, err := New()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET must be at least 32 bytes")
}

func TestNew_JWTSecretExactMinLength(t *testing.T) {
	setRequiredEnv(t) // validSecret is exactly 32 bytes

	cfg, err := New()

	require.NoError(t, err)
	assert.Equal(t, validSecret, cfg.JWTSecret)
}

func TestNew_DBPoolDefaults(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := New()

	require.NoError(t, err)
	assert.Equal(t, 25, cfg.DBMaxOpenConns)
	assert.Equal(t, 5, cfg.DBMaxIdleConns)
	assert.Equal(t, 300, cfg.DBConnMaxLifetimeSecs)
}

func TestNew_DBPoolFromEnv(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_MAX_IDLE_CONNS", "10")
	t.Setenv("DB_CONN_MAX_LIFETIME_SECS", "600")

	cfg, err := New()

	require.NoError(t, err)
	assert.Equal(t, 50, cfg.DBMaxOpenConns)
	assert.Equal(t, 10, cfg.DBMaxIdleConns)
	assert.Equal(t, 600, cfg.DBConnMaxLifetimeSecs)
}

func TestNew_JWTExpiryFromEnv(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("JWT_EXPIRY_HOURS", "48")

	cfg, err := New()

	require.NoError(t, err)
	assert.Equal(t, 48, cfg.JWTExpiryHours)
}

func TestNew_AutoMigrateDefaultFalse(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := New()

	require.NoError(t, err)
	assert.False(t, cfg.AutoMigrate)
}

func TestNew_AutoMigrateFromEnv(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("DB_AUTO_MIGRATE", "true")

	cfg, err := New()

	require.NoError(t, err)
	assert.True(t, cfg.AutoMigrate)
}
