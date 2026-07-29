package persistence

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	DatabaseURL    string
	MaxConns       int32
	MinConns       int32
	ConnectTimeout time.Duration
	SearchPath     string
}

func PostgresConfigFromEnv() (PostgresConfig, error) {
	cfg := PostgresConfig{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		MaxConns:       10,
		MinConns:       1,
		ConnectTimeout: 10 * time.Second,
		SearchPath:     os.Getenv("POSTGRES_SEARCH_PATH"),
	}
	if cfg.DatabaseURL == "" {
		return PostgresConfig{}, fmt.Errorf("DATABASE_URL is required")
	}
	var err error
	if raw := os.Getenv("POSTGRES_MAX_CONNS"); raw != "" {
		cfg.MaxConns, err = parsePositiveInt32("POSTGRES_MAX_CONNS", raw)
		if err != nil {
			return PostgresConfig{}, err
		}
	}
	if raw := os.Getenv("POSTGRES_MIN_CONNS"); raw != "" {
		value, parseErr := strconv.ParseInt(raw, 10, 32)
		if parseErr != nil || value < 0 {
			return PostgresConfig{}, fmt.Errorf("POSTGRES_MIN_CONNS must be a non-negative integer")
		}
		cfg.MinConns = int32(value)
	}
	if cfg.MinConns > cfg.MaxConns {
		return PostgresConfig{}, fmt.Errorf("POSTGRES_MIN_CONNS cannot exceed POSTGRES_MAX_CONNS")
	}
	if raw := os.Getenv("POSTGRES_CONNECT_TIMEOUT"); raw != "" {
		cfg.ConnectTimeout, err = time.ParseDuration(raw)
		if err != nil || cfg.ConnectTimeout <= 0 {
			return PostgresConfig{}, fmt.Errorf("POSTGRES_CONNECT_TIMEOUT must be a positive duration")
		}
	}
	return cfg, nil
}

func OpenPostgres(ctx context.Context, cfg PostgresConfig) (*pgxpool.Pool, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse PostgreSQL configuration: %w", err)
	}
	if cfg.MaxConns > 0 {
		poolConfig.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns >= 0 {
		poolConfig.MinConns = cfg.MinConns
	}
	if cfg.ConnectTimeout > 0 {
		poolConfig.ConnConfig.ConnectTimeout = cfg.ConnectTimeout
	}
	if cfg.SearchPath != "" {
		poolConfig.ConnConfig.RuntimeParams["search_path"] = cfg.SearchPath
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}
	return pool, nil
}

func parsePositiveInt32(name, raw string) (int32, error) {
	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return int32(value), nil
}
