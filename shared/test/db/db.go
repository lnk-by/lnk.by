package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5"
	"github.com/lnk.by/shared/db"
	"github.com/stretchr/testify/assert"
)

var conn *pgx.Conn

func panicErr(message string, err error) {
	panic(fmt.Errorf("%s: %w", message, err))
}

func StartDb(ctx context.Context, dbUrl string, username string, password string, dbDir string) func() {
	tempDir, err := os.MkdirTemp("", "pgdata-*")
	if err != nil {
		panicErr("failed to create temp dir", err)
	}

	// Start embedded Postgres
	postgres := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Username(username).
			Password(password).
			RuntimePath(tempDir).
			Port(9876),
	)

	if err := postgres.Start(); err != nil {
		panicErr("failed to start embedded Postgres", err)
	}

	// Connect using pgx
	conn, err = pgx.Connect(ctx, dbUrl)
	if err != nil {
		if stopErr := postgres.Stop(); stopErr != nil {
			err = errors.Join(err, stopErr)
		}
		panicErr("failed to connect to DB", err)
	}

	// Run create.sql
	if err := runSQLFile(ctx, conn, dbDir+"/create.sql"); err != nil {
		if stopErr := postgres.Stop(); stopErr != nil {
			err = errors.Join(err, stopErr)
		}
		panicErr("failed to run create.sql", err)
	}

	if err := db.Init(ctx, dbUrl, "test", "test"); err != nil {
		panicErr("failed to initialize DB", err)
	}

	return func() {
		// Run drop.sql for cleanup
		if err := runSQLFile(ctx, conn, dbDir+"/drop.sql"); err != nil {
			fmt.Printf("Warning: Failed to run drop.sql: %v\n", err)
		}

		_ = conn.Close(ctx)
		_ = postgres.Stop()
		_ = os.RemoveAll(tempDir)
	}
}

func runSQLFile(ctx context.Context, conn *pgx.Conn, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, string(content))
	return err
}

func truncateTable(t *testing.T, table string) {
	_, err := conn.Exec(t.Context(), fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
	assert.NoError(t, err)
}

func WithTable(t *testing.T, table string, f func()) {
	truncateTable(t, table)
	defer truncateTable(t, table)

	f()
}
