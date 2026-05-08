package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"hello_world/internal/config"
	"hello_world/internal/model"
)

// Service defines the contract for interacting with the database layer.
type Service interface {
Health(ctx context.Context) map[string]string
Close() error
GetDB() *gorm.DB
}

type service struct {
	db *gorm.DB
}

// New opens a Postgres connection, configures the connection pool, verifies
// connectivity, and runs auto-migrations. It returns an error instead of
// calling log.Fatal so the caller can handle the failure gracefully.
func New(cfg *config.Config) (Service, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s search_path=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBSchema,
	)

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting underlying sql.DB: %w", err)
	}

	// Connection pool tuning.
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetimeSecs) * time.Second)
	// Keep idle timeout proportional to max lifetime while preserving previous 2m default behavior.
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.DBConnMaxLifetimeSecs/2) * time.Second)

	// Fail fast if the database is unreachable at startup.
	if err := sqlDB.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	// Auto-migrate is opt-in to avoid schema drift in production.
	if cfg.AutoMigrate {
		if err := db.AutoMigrate(model.Registry...); err != nil {
			return nil, fmt.Errorf("running auto-migrations: %w", err)
		}
	}

	return &service{db: db}, nil
}

func (s *service) GetDB() *gorm.DB {
	return s.db
}

func (s *service) Health(ctx context.Context) map[string]string {
	sqlDB, err := s.db.DB()
	if err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}
	stats := sqlDB.Stats()
	return map[string]string{
		"status":       "up",
		"open_conns":   fmt.Sprintf("%d", stats.OpenConnections),
		"in_use_conns": fmt.Sprintf("%d", stats.InUse),
		"idle_conns":   fmt.Sprintf("%d", stats.Idle),
	}
}

func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("getting underlying sql.DB for close: %w", err)
	}
	return sqlDB.Close()
}
