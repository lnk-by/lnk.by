package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5"
	"github.com/lnk.by/shared/db"
	"github.com/lnk.by/shared/service"
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

func truncateTable(ctx context.Context, table string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
	return err
}

func TruncateTable(t *testing.T, table string) {
	err := truncateTable(t.Context(), table)
	assert.NoError(t, err)

	t.Cleanup(func() {
		// the t.Context() is already cancelled at this point
		err := truncateTable(context.Background(), table)
		assert.NoError(t, err)
	})
}

func Create[T service.Creatable](t *testing.T, createSQL service.CreateSQL[T], entity T) T {
	bytes, err := json.Marshal(entity)
	assert.NoError(t, err)

	status, body := service.Create(t.Context(), createSQL, bytes)
	assert.Equal(t, http.StatusCreated, status)

	var created T
	err = json.Unmarshal([]byte(body), &created)
	assert.NoError(t, err)

	return created
}

func Retrieve[T service.FieldsPtrsAware](t *testing.T, retrieveSQL service.RetrieveSQL[T], id string) T {
	status, body := service.Retrieve(t.Context(), retrieveSQL, id)
	assert.Equal(t, http.StatusOK, status)

	var retrieved T
	err := json.Unmarshal([]byte(body), &retrieved)
	assert.NoError(t, err)

	return retrieved
}

func List[T service.FieldsPtrsAware](t *testing.T, listSQL service.ListSQL[T], offset int, limit int) []T {
	status, body := service.List(t.Context(), listSQL, offset, limit)
	assert.Equal(t, http.StatusOK, status)

	var listed []T
	err := json.Unmarshal([]byte(body), &listed)
	assert.NoError(t, err)

	return listed
}
