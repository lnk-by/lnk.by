package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitFromEnvironment(ctx context.Context) error {
	slog.Info("Connecting to database ", "url", os.Getenv("DB_URL"), "user", os.Getenv("DB_USER"), "password", os.Getenv("DB_PASSWORD"))
	if err := Init(ctx, os.Getenv("DB_URL"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")); err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	return nil
}

const maxPingRetries = 10

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

	// retry pings if Postgres is not ready yet -- can happen under docker compose
	for retries := 0; retries < maxPingRetries; retries++ {
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

func RunScript(ctx context.Context, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}
	queries := strings.Split(string(content), ";")

	conn, err := Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get DB connection: %w", err)
	}

	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" || strings.HasPrefix(q, "--") {
			continue
		}

		if _, err := conn.Exec(ctx, q); err != nil {
			return fmt.Errorf("failed to execute query: %w\nSQL: %s", err, q)
		}

		slog.Info("Executed", "sql", q)
	}

	return nil
}
