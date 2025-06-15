package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var pool *pgxpool.Pool

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

	//if err := pool.Ping(ctx); err != nil {
	//	return fmt.Errorf("failed to ping DB pool: %w", err)
	//}

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
