package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/seeds"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"log"
	"log/slog"
	"os"
	"time"
)

var DB *sql.DB

func NewPostgresConnection(cfg config.Config) *sql.DB {

	var dsn string

	if cfg.DatabaseConfig.ConnString != "" {
		slog.Info("Using custom PostgreSQL connection string")
		dsn = cfg.DatabaseConfig.ConnString
	} else {
		slog.Info("Using PostgreSQL connection parameters")
		if cfg.DatabaseConfig.Host == "" || cfg.DatabaseConfig.Port == "" || cfg.DatabaseConfig.User == "" || cfg.DatabaseConfig.Password == "" || cfg.DatabaseConfig.Name == "" {
			slog.Error("Missing required PostgreSQL connection parameters")
			log.Fatal("Error: Missing required PostgreSQL connection parameters")
		}
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.User, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.Name,
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Failed to open PostgreSQL connection", "error", err)
		log.Fatalf("Error opening PostgreSQL connection: %v", err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

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
	DB = db
	return db
}
