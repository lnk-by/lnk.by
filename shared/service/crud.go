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

type validatable interface {
	Validate() error
}

type idAware interface {
	WithId(id string)
}

type FieldsValsAware interface {
	validatable
	idAware
	FieldsVals() []any
}

func decode[T any](content io.ReadCloser) (t T, err error) {
	defer func() {
		if closeErr := content.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	t = inst[T]()
	if err = json.NewDecoder(content).Decode(t); err != nil {
		err = fmt.Errorf("failed to decode %T from JSON: %w", t, err)
	}
	return
}

func Create[T FieldsValsAware](ctx context.Context, createSQL CreateSQL[T], content io.ReadCloser) (int, string) {
	status, body := CreateWithRetries(ctx, createSQL, content, 1)
	return status, body
}

func CreateWithRetries[T FieldsValsAware](ctx context.Context, createSQL CreateSQL[T], content io.ReadCloser, maxAttempts int) (int, string) {
	t, err := decode[T](content)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create %T: %w", t, err).Error()
	}

	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		for i := 0; i < maxAttempts; i++ {
			if err := t.Validate(); err != nil {
				return failed(http.StatusBadRequest, fmt.Errorf("failed to validate %v: %w", t, err))
			}

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

func Update[T FieldsValsAware](ctx context.Context, updateSQL UpdateSQL[T], id string, content io.ReadCloser) (int, string) {
	t, err := decode[T](content)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to update %T: %w", t, err).Error()
	}

	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		if err := t.Validate(); err != nil {
			return failed(http.StatusBadRequest, fmt.Errorf("failed to validate %v: %w", t, err))
		}
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
