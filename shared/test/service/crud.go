package service

import (
	"encoding/json"
	"github.com/lnk.by/shared/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

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
