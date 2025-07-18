package service

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/lnk.by/shared/service"
	"github.com/stretchr/testify/assert"
)

func Create[T service.Creatable](t *testing.T, createSQL service.CreateSQL[T], entity T) T {
	status, body := service.Create(t.Context(), createSQL, marshal(t, entity))
	assert.Equal(t, http.StatusCreated, status)

	return unmarshal[T](t, body)
}

func Retrieve[K any, T service.Retrievable[K]](t *testing.T, retrieveSQL service.RetrieveSQL[T], id string) T {
	status, body := service.Retrieve(t.Context(), retrieveSQL, id)
	assert.Equal(t, http.StatusOK, status)

	return unmarshal[T](t, body)
}

func Update[K any, T service.Updatable[K]](t *testing.T, updateSQL service.UpdateSQL[T], id string, entity T) T {
	status, body := service.Update(t.Context(), updateSQL, id, marshal(t, entity))
	assert.Equal(t, http.StatusOK, status)

	return unmarshal[T](t, body)
}

func Delete[K any, T service.Identifiable[K]](t *testing.T, deleteSQL service.DeleteSQL[T], id string) {
	status, body := service.Delete(t.Context(), deleteSQL, id)
	assert.Equal(t, http.StatusNoContent, status)
	assert.Len(t, body, 0)
}

func List[K any, T service.Retrievable[K]](t *testing.T, listSQL service.ListSQL[T], offset int, limit int) []T {
	status, body := service.List(t.Context(), listSQL, offset, limit)
	assert.Equal(t, http.StatusOK, status)

	return unmarshal[[]T](t, body)
}

func marshal(t *testing.T, entity any) []byte {
	bytes, err := json.Marshal(entity)
	assert.NoError(t, err)

	return bytes
}

func unmarshal[T any](t *testing.T, body string) T {
	var res T
	err := json.Unmarshal([]byte(body), &res)
	assert.NoError(t, err)

	return res
}
