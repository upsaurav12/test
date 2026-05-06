package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setRequiredEnv sets the minimum env vars needed for a valid Config.
func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("GONE_DB_DATABASE", "testdb")
	t.Setenv("GONE_DB_USERNAME", "testuser")
	t.Setenv("GONE_DB_PASSWORD", "testpass")
	t.Setenv("JWT_SECRET", "supersecretkey")
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
