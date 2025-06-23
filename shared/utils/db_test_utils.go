package utils

import (
	"context"
	"fmt"
	"os"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5"
	"github.com/lnk.by/shared/db"
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
