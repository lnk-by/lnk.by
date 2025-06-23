package customer

import (
	"context"
	"net/http"
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
	adam := Customer{Email: "adam@human.net", Name: "Adam"}
	adamBytes, err := json.Marshal(adam)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	status, body := service.Create(context.Background(), CreateSQL, adamBytes)
	assert.Equal(t, http.StatusCreated, status)
	var created Customer
	if err := json.Unmarshal([]byte(body), &created); err != nil {
		assert.Fail(t, err.Error())
	}

	status, body = service.Retrieve(context.Background(), RetrieveSQL, created.ID)
	assert.Equal(t, http.StatusOK, status)
	var retrieved Customer
	if err := json.Unmarshal([]byte(body), &retrieved); err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, created, retrieved)

	status, body = service.List(context.Background(), ListSQL, 0, 10)
	assert.Equal(t, http.StatusOK, status)

	var listed []Customer
	if err := json.Unmarshal([]byte(body), &listed); err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, 1, len(listed))
	assert.Equal(t, created, listed[0])
}
