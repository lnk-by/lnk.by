package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lnk.by/shared/db"
)

const IdParam = "id"

type CreateSQL[T any] string
type RetrieveSQL[T any] string
type UpdateSQL[T any] string
type DeleteSQL[T any] string
type ListSQL[T any] string

type failure struct {
	Error string `json:"error"`
}

func failed(status int, err error) (int, string) {
	jsonBytes, err := json.Marshal(failure{Error: err.Error()})
	if err != nil {
		return status, http.StatusText(status)
	}

	return status, string(jsonBytes)
}

func withConn[T any](ctx context.Context, f func(conn *pgxpool.Conn) (int, T, error)) (int, T, error) {
	conn, err := db.Get(ctx)
	if err != nil {
		var zero T
		return http.StatusInternalServerError, zero, fmt.Errorf("failed to get DB connection: %w", err)
	}
	defer conn.Release()

	return f(conn)
}

type FieldsValsAware interface {
	Validate() error
	FieldsVals() []any
}

type Creatable interface {
	FieldsValsAware
	Generate()
}

type retriable interface {
	MaxAttempts() int
}

func CreateFromReqBody[T Creatable](ctx context.Context, createSQL CreateSQL[T], body io.ReadCloser) (int, string) {
	content, err := io.ReadAll(body)
	if err != nil {
		return failed(http.StatusInternalServerError, fmt.Errorf("failed to read request body: %w", err))
	}
	return Create(ctx, createSQL, content)
}

func Create[T Creatable](ctx context.Context, createSQL CreateSQL[T], content []byte) (int, string) {
	t := inst[T]()
	if err := json.Unmarshal(content, t); err != nil {
		return failed(http.StatusBadRequest, fmt.Errorf("failed to unmarshal %T from JSON: %w", t, err))
	}
	if err := t.Validate(); err != nil {
		return failed(http.StatusBadRequest, fmt.Errorf("failed to validate %T: %w", t, err))
	}
	return CreateRecord(ctx, createSQL, t, 0)
}

func CreateRecord[T Creatable](ctx context.Context, createSQL CreateSQL[T], t T, generateFromIteration int) (int, string) {
	maxAttempts := 1
	if t, ok := any(t).(retriable); ok {
		maxAttempts = t.MaxAttempts()
	}

	return marshal(withConn(ctx, func(conn *pgxpool.Conn) (int, T, error) {
		for i := 0; i < maxAttempts; i++ {
			if i >= generateFromIteration {
				t.Generate()
			}
			if _, err := conn.Exec(ctx, string(createSQL), t.FieldsVals()...); err != nil {
				if isDuplicateKeyError(err) {
					continue // try again
				}
				return http.StatusInternalServerError, t, fmt.Errorf("failed to insert %T %v: %w", t, t, err)
			}

			return http.StatusCreated, t, nil
		}

		return http.StatusConflict, t, fmt.Errorf("failed to create unique identifier for %T", t)
	}))
}

func isDuplicateKeyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}

type FieldsPtrsAware interface {
	FieldsPtrs() []any
}

func inst[T any]() T {
	return reflect.New(reflect.TypeFor[T]().Elem()).Interface().(T)
}

func Retrieve[T FieldsPtrsAware](ctx context.Context, retrieveSQL RetrieveSQL[T], id string) (int, string) {
	return marshal(retrieve(ctx, retrieveSQL, id))
}

func marshal[T any](status int, t T, err error) (int, string) {
	if err != nil {
		return failed(status, err)
	}

	var jsonBytes []byte
	if !reflect.ValueOf(t).IsNil() {
		jsonBytes, err = json.Marshal(t)
		if err != nil {
			return failed(http.StatusInternalServerError, fmt.Errorf("failed to marshal the %T %v: %w", t, t, err))
		}
	}
	return status, string(jsonBytes)
}

func RetrieveValueAndMarshalError[T FieldsPtrsAware](ctx context.Context, retrieveSQL RetrieveSQL[T], id string) (int, T, string) {
	status, value, err := retrieve(ctx, retrieveSQL, id)
	if err != nil {
		_, strErr := failed(status, err)
		return status, value, strErr
	}

	return status, value, ""
}

func retrieve[T FieldsPtrsAware](ctx context.Context, retrieveSQL RetrieveSQL[T], id string) (int, T, error) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, T, error) {
		t := inst[T]()
		if err := conn.QueryRow(ctx, string(retrieveSQL), id).Scan(t.FieldsPtrs()...); err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				return http.StatusNotFound, t, fmt.Errorf("failed to retrieve the %T with id %q: %w", t, id, err)
			case errors.Is(err, pgx.ErrTooManyRows):
				return http.StatusConflict, t, fmt.Errorf("failed to retrieve the %T with id %q: %w", t, id, err)
			default:
				return http.StatusInternalServerError, t, fmt.Errorf("failed to retrieve the %T with id %q: %w", t, id, err)
			}
		}

		return http.StatusOK, t, nil
	})
}

type Updatable interface {
	FieldsValsAware
	WithId(id string)
}

func UpdateFromReqBody[T Updatable](ctx context.Context, updateSQL UpdateSQL[T], id string, body io.ReadCloser) (int, string) {
	content, err := io.ReadAll(body)
	if err != nil {
		return failed(http.StatusInternalServerError, fmt.Errorf("failed to read request body: %w", err))
	}
	return Update(ctx, updateSQL, id, content)
}

func Update[T Updatable](ctx context.Context, updateSQL UpdateSQL[T], id string, content []byte) (int, string) {
	t := inst[T]()
	if err := json.Unmarshal(content, t); err != nil {
		return failed(http.StatusBadRequest, fmt.Errorf("failed to unmarshal %T from JSON: %w", t, err))
	}
	if err := t.Validate(); err != nil {
		return failed(http.StatusBadRequest, fmt.Errorf("failed to validate %T: %w", t, err))
	}

	return marshal(withConn(ctx, func(conn *pgxpool.Conn) (int, T, error) {
		t.WithId(id)
		commandTag, err := conn.Exec(ctx, string(updateSQL), t.FieldsVals()...)
		switch {
		case err != nil:
			return http.StatusInternalServerError, t, fmt.Errorf("failed to update %T %v: %w", t, t, err)
		case commandTag.RowsAffected() == 0:
			return http.StatusNotFound, t, fmt.Errorf("failed to update %T %v: %w", t, t, pgx.ErrNoRows)
		case commandTag.RowsAffected() > 1: // TODO is it too late?
			return http.StatusConflict, t, fmt.Errorf("failed to update %T %v: %w", t, t, pgx.ErrTooManyRows)
		}

		return http.StatusOK, t, nil
	}))
}

func Delete[T any](ctx context.Context, deleteSQL DeleteSQL[T], id string) (int, string) {
	return marshal(withConn(ctx, func(conn *pgxpool.Conn) (int, T, error) {
		commandTag, err := conn.Exec(ctx, string(deleteSQL), id)
		var t T // only to build the error
		switch {
		case err != nil:
			return http.StatusInternalServerError, t, fmt.Errorf("failed to delete %T with id %v: %w", t, id, err)
		case commandTag.RowsAffected() == 0:
			return http.StatusNotFound, t, fmt.Errorf("failed to delete %T with id %v: %w", t, id, pgx.ErrNoRows)
		case commandTag.RowsAffected() > 1: // TODO is it too late?
			return http.StatusNotFound, t, fmt.Errorf("failed to delete %T with id %v: %w", t, id, pgx.ErrTooManyRows)
		}

		return http.StatusNoContent, t, nil
	}))
}

func List[T FieldsPtrsAware](ctx context.Context, listSQL ListSQL[T], offset int, limit int) (int, string) {
	return marshal(withConn[[]T](ctx, func(conn *pgxpool.Conn) (int, []T, error) {
		sql := string(listSQL)

		rows, err := conn.Query(ctx, sql, offset, limit)
		if err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("failed to execute list query for %T: %w", new(T), err)
		}
		defer rows.Close()

		ts := make([]T, 0)
		for rows.Next() {
			t := inst[T]()
			if err := rows.Scan(t.FieldsPtrs()...); err != nil {
				return http.StatusInternalServerError, ts, fmt.Errorf("failed to scan row to build the %T: %w", t, err)
			}

			ts = append(ts, t)
		}

		return http.StatusOK, ts, nil
	}))
}

func UUID() string {
	return uuid.Must(uuid.NewV1()).String()
}

var (
	ErrNameRequired      = errors.New("name is required")
	ErrIDManagedByServer = errors.New("ID is managed by the server")
)
