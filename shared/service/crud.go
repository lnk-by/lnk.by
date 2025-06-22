package service

import (
	"bytes"
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

func succeeded(status int, t any) (int, string) {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return failed(http.StatusInternalServerError, fmt.Errorf("failed to marshal the %T %v: %w", t, t, err))
	}

	return status, string(jsonBytes)
}

func withConn(ctx context.Context, f func(conn *pgxpool.Conn) (int, string)) (int, string) {
	conn, err := db.Get(ctx)
	if err != nil {
		return failed(http.StatusInternalServerError, fmt.Errorf("failed to get DB connection: %w", err))
	}
	defer conn.Release()

	return f(conn)
}

type FieldsValsAware interface {
	Validate() error
	FieldsVals() []any
}

func decode[T FieldsValsAware](content []byte) (t T, status int, err error) {
	t = inst[T]()
	if err = json.Unmarshal(content, t); err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("failed to unmarshal %T from JSON: %w", t, err)
	}

	if err = t.Validate(); err != nil {
		status = http.StatusBadRequest
		err = fmt.Errorf("failed to validate %T %v: %w", t, t, err)
	}

	return
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
	t, status, err := decode[T](content)
	if err != nil {
		return failed(status, fmt.Errorf("failed to create %T: %w", t, err))
	}

	maxAttempts := 1
	if t, ok := any(t).(retriable); ok {
		maxAttempts = t.MaxAttempts()
	}

	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		for i := 0; i < maxAttempts; i++ {
			t.Generate()
			if _, err := conn.Exec(ctx, string(createSQL), t.FieldsVals()...); err != nil {
				if isDuplicateKeyError(err) {
					continue // try again
				}
				return failed(http.StatusInternalServerError, fmt.Errorf("failed to insert %T %v: %w", t, t, err))
			}

			return succeeded(http.StatusCreated, t)
		}

		return http.StatusConflict, fmt.Errorf("failed to create unique identifier").Error()
	})
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
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		t := inst[T]()
		if err := conn.QueryRow(ctx, string(retrieveSQL), id).Scan(t.FieldsPtrs()...); err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				return failed(http.StatusNotFound, fmt.Errorf("failed to retrieve the %T with id %q: %w", t, id, err))
			case errors.Is(err, pgx.ErrTooManyRows):
				return failed(http.StatusConflict, fmt.Errorf("failed to retrieve the %T with id %q: %w", t, id, err))
			default:
				return failed(http.StatusInternalServerError, fmt.Errorf("failed to retrieve the %T with id %q: %w", t, id, err))
			}
		}

		return succeeded(http.StatusOK, t)
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
	t, status, err := decode[T](content)
	if err != nil {
		return failed(status, fmt.Errorf("failed to update %T: %w", t, err))
	}

	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		t.WithId(id)
		commandTag, err := conn.Exec(ctx, string(updateSQL), t.FieldsVals()...)
		switch {
		case err != nil:
			return failed(http.StatusInternalServerError, fmt.Errorf("failed to update %T %v: %w", t, t, err))
		case commandTag.RowsAffected() == 0:
			return failed(http.StatusNotFound, fmt.Errorf("failed to update %T %v: %w", t, t, pgx.ErrNoRows))
		case commandTag.RowsAffected() > 1: // TODO is it too late?
			return failed(http.StatusConflict, fmt.Errorf("failed to update %T %v: %w", t, t, pgx.ErrTooManyRows))
		}

		return succeeded(http.StatusOK, t)
	})
}

func Delete[T any](ctx context.Context, deleteSQL DeleteSQL[T], id string) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		commandTag, err := conn.Exec(ctx, string(deleteSQL), id)
		switch {
		case err != nil:
			var t T // only to build the error
			return failed(http.StatusInternalServerError, fmt.Errorf("failed to delete %T with id %v: %w", t, id, err))
		case commandTag.RowsAffected() == 0:
			var t T // only to build the error
			return failed(http.StatusNotFound, fmt.Errorf("failed to delete %T with id %v: %w", t, id, pgx.ErrNoRows))
		case commandTag.RowsAffected() > 1: // TODO is it too late?
			var t T // only to build the error
			return failed(http.StatusNotFound, fmt.Errorf("failed to delete %T with id %v: %w", t, id, pgx.ErrTooManyRows))
		}

		return http.StatusNoContent, ""
	})
}

func List[T FieldsPtrsAware](ctx context.Context, listSQL ListSQL[T], offset int, limit int) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		sql := string(listSQL)

		rows, err := conn.Query(ctx, sql, offset, limit)
		if err != nil {
			return failed(http.StatusInternalServerError, fmt.Errorf("failed to execute list query for %T: %w", new(T), err))
		}
		defer rows.Close()

		var buf bytes.Buffer

		buf.WriteByte('[')

		t := inst[T]()
		if _, err := pgx.ForEachRow(rows, t.FieldsPtrs(), func() error {
			if buf.Len() > 1 { // the buf contains at least one marshalled row
				buf.WriteByte(',')
			}

			jsonBytes, err := json.Marshal(t)
			if err != nil {
				return fmt.Errorf("failed to marshal entity: %w", err)
			}
			buf.Write(jsonBytes)

			return nil
		}); err != nil {
			return failed(http.StatusInternalServerError, fmt.Errorf("failed to iterate over rows: %w", err))
		}

		buf.WriteByte(']')

		return http.StatusOK, buf.String()
	})
}

func UUID() string {
	return uuid.Must(uuid.NewV1()).String()
}

var (
	ErrNameRequired      = errors.New("name is required")
	ErrIDManagedByServer = errors.New("ID is managed by the server")
)
