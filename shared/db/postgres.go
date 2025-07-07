package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitFromEnvironment(ctx context.Context) error {
	slog.Info("Connecting to database ", "url", os.Getenv("DB_URL"), "user", os.Getenv("DB_USER"), "password", os.Getenv("DB_PASSWORD"))
	if err := Init(ctx, os.Getenv("DB_URL"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}
	return nil
}

const maxRetries = 10

func Init(ctx context.Context, dbUrl string, user string, password string) error {
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return fmt.Errorf("failed to parse DB config %q: %w", dbUrl, err)
	}

	config.ConnConfig.User = user
	config.ConnConfig.Password = password

	config.MaxConns = 5
	config.MinConns = 1
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 30 * time.Second

	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to build DB pool: %w", err)
	}

	for retries := 0; retries < maxRetries; retries++ {
		if err = pool.Ping(ctx); err != nil {
			slog.Info("Waiting for DB pool", "error", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	if err != nil {
		return fmt.Errorf("failed to ping DB pool: %w", err)
	}

	return nil
}

func Get(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire DB connection: %w", err)
	}

	//if err := conn.Ping(ctx); err != nil {
	//	return nil, fmt.Errorf("failed to ping DB connection: %w", err)
	//}

	return conn, nil
}
