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
	utils.StartDb("postgres://test:test@localhost:9876/postgres?sslmode=disable", "test", "test", "../../db")

	// Run tests
	code := m.Run()

	utils.StopDb("../../db")

	os.Exit(code)
}

func TestListEmpty(t *testing.T) {
	utils.CleanupTestDatabase(t, "customer")
	status, body := service.List(context.Background(), ListSQL, 0, 10)
	assert.Equal(t, 200, status)
	assert.Equal(t, "[]", body)

	var customers []Customer
	if err := json.Unmarshal([]byte(body), &customers); err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, 0, len(customers))
}

func TestCreateAndGet(t *testing.T) {
	utils.CleanupTestDatabase(t, "customer")
	adam := Customer{Email: "adam@human.net", Name: "Adam"}
	created := utils.Create(t, CreateSQL, &adam)

	retrieved := utils.Retrieve(t, RetrieveSQL, created.ID)
	assert.Equal(t, created, retrieved)

	listed := utils.List(t, ListSQL, 0, 10)
	assert.Equal(t, 1, len(listed))
	assert.Equal(t, created, listed[0])
}
