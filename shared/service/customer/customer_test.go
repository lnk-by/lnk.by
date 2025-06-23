package customer

import (
	"context"
	"os"
	"testing"

	"encoding/json"

	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(
		func() int {
			ctx := context.Background()
			utils.StartDb(ctx, "postgres://test:test@localhost:9876/postgres?sslmode=disable", "test", "test", "../../db")
			defer utils.StopDb(ctx, "../../db")

			return m.Run() // Run tests
		}(),
	)
}

func TestListEmpty(t *testing.T) {
	utils.TruncateTable(t, "customer")
	status, body := service.List(t.Context(), ListSQL, 0, 10)
	assert.Equal(t, 200, status)
	assert.Equal(t, "[]", body)

	var customers []Customer
	err := json.Unmarshal([]byte(body), &customers)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(customers))
}

func TestCreateAndGet(t *testing.T) {
	utils.TruncateTable(t, "customer")
	adam := Customer{Email: "adam@human.net", Name: "Adam"}
	created := utils.Create(t, CreateSQL, &adam)

	retrieved := utils.Retrieve(t, RetrieveSQL, created.ID)
	assert.Equal(t, created, retrieved)

	listed := utils.List(t, ListSQL, 0, 10)
	assert.Equal(t, 1, len(listed))
	assert.Equal(t, created, listed[0])
}
