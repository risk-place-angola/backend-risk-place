package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/seeds"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
)

const (
	// Database connection pool settings for high-load production
	// Rule of thumb: maxOpenConns = ((core_count * 2) + effective_spindle_count)
	// For cloud deployments: start with 100-200 for high concurrency
	maxOpenConns    = 100 // Maximum number of open connections to the database
	maxIdleConns    = 25  // Keep 25% idle for quick reuse
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 5 * time.Minute // Close idle connections after 5 minutes
)

func NewPostgresConnection(cfg config.Config) *sql.DB {
	db, err := Connect(cfg)
	if err != nil {
		slog.Error("Failed to connect to PostgreSQL database", "error", err)
		log.Fatalf("Error connecting to PostgreSQL database: %v", err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	slog.Info("Database connection pool configured",
		"maxOpenConns", maxOpenConns,
		"maxIdleConns", maxIdleConns,
		"connMaxLifetime", connMaxLifetime,
		"connMaxIdleTime", connMaxIdleTime,
	)

	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping PostgreSQL database", "error", err)
		log.Fatalf("Error pinging PostgreSQL database: %v", err)
		os.Exit(1)
	}

	// seeds
	if cfg.IsDevelopment() {
		slog.Info("Running database seeds for development environment")
		if err := seeds.RunAll(context.Background(), db); err != nil {
			slog.Error("database seeds failed", "err", err)
		}
		slog.Info("Database seeds executed successfully")
	}

	slog.Info("PostgreSQL connection established successfully")
	return db
}

func ClosePostgresConnection(db *sql.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			slog.Error("Failed to close PostgreSQL connection", "error", err)
			log.Fatalf("Error closing PostgreSQL connection: %v", err)
		} else {
			slog.Info("PostgreSQL connection closed successfully")
		}
	}
}

func Connect(cfg config.Config) (*sql.DB, error) {
	dsn := GetDBConnectionString(cfg)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDBConnectionString(cfg config.Config) string {
	if cfg.DatabaseConfig.ConnString != "" {
		return cfg.DatabaseConfig.ConnString
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.User, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.Name,
	)
}
