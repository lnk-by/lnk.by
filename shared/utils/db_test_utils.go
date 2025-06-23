package utils

import (
	"context"
	"encoding/json"
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

var postgres *embeddedpostgres.EmbeddedPostgres
var tempDir string
var conn *pgx.Conn

func StartDb(dbUrl string, username string, password string, dbDir string) {
	// Create temp directory
	var err error
	tempDir, err = os.MkdirTemp("", "pgdata-*")
	if err != nil {
		panic(fmt.Sprintf("Failed to create temp dir: %v", err))
	}

	// Start embedded Postgres
	postgres = embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Username(username).
			Password(password).
			RuntimePath(tempDir).
			Port(9876),
	)

	if err := postgres.Start(); err != nil {
		panic(fmt.Sprintf("Failed to start embedded Postgres: %v", err))
	}

	// Connect using pgx
	conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		postgres.Stop()
		panic(fmt.Sprintf("Failed to connect to DB: %v", err))
	}

	// Run create.sql
	if err := runSQLFile(conn, dbDir+"/create.sql"); err != nil {
		postgres.Stop()
		panic(fmt.Sprintf("Failed to run create.sql: %v", err))
	}

	if err := db.Init(context.Background(), dbUrl, "test", "test"); err != nil {
		os.Exit(1)
	}
}

func StopDb(dbDir string) {
	// Run drop.sql for cleanup
	if err := runSQLFile(conn, dbDir+"/db/drop.sql"); err != nil {
		fmt.Printf("Warning: Failed to run drop.sql: %v\n", err)
	}

	_ = conn.Close(context.Background())
	_ = postgres.Stop()
	_ = os.RemoveAll(tempDir)

}

func runSQLFile(conn *pgx.Conn, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), string(content))
	return err
}

func ClearDatabase(table string) error {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
	return err
}

func CleanupTestDatabase(t *testing.T, table string) {
	ClearDatabase(table)
	t.Cleanup(func() {
		if err := ClearDatabase(table); err != nil {
			t.Errorf("failed to clear DB after test: %v", err)
		}
	})
}

func Create[T service.Creatable](t *testing.T, createSQL service.CreateSQL[T], entity T) T {
	bytes, err := json.Marshal(entity)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	status, body := service.Create(context.Background(), createSQL, bytes)
	assert.Equal(t, http.StatusCreated, status)
	var created T
	if err := json.Unmarshal([]byte(body), &created); err != nil {
		assert.Fail(t, err.Error())
	}
	return created
}

func Retrieve[T service.FieldsPtrsAware](t *testing.T, retrieveSQL service.RetrieveSQL[T], id string) T {
	status, body := service.Retrieve(context.Background(), retrieveSQL, id)
	assert.Equal(t, http.StatusOK, status)
	var retrieved T
	if err := json.Unmarshal([]byte(body), &retrieved); err != nil {
		assert.Fail(t, err.Error())
	}
	return retrieved
}

func List[T service.FieldsPtrsAware](t *testing.T, listSQL service.ListSQL[T], offset int, limit int) []T {
	status, body := service.List(context.Background(), listSQL, offset, limit)
	assert.Equal(t, http.StatusOK, status)

	var listed []T
	if err := json.Unmarshal([]byte(body), &listed); err != nil {
		assert.Fail(t, err.Error())
	}
	return listed
}
