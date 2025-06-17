package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lnk.by/shared/db"
)

type CreateSQL[T any] string
type RetrieveSQL[T any] string
type UpdateSQL[T any] string
type DeleteSQL[T any] string
type ListSQL[T any] string

func withConn(ctx context.Context, f func(conn *pgxpool.Conn) (int, string)) (int, string) {
	conn, err := db.Get(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get DB connection: %w", err).Error()
	}
	defer conn.Release()

	return f(conn)
}

type validatable interface {
	Validate() error
}

type fieldsValsAware interface {
	validatable
	FieldsVals() []any
}

func Create[T fieldsValsAware](ctx context.Context, createSQL CreateSQL[T], t T) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		if err := t.Validate(); err != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to validate %v: %w", t, err).Error()
		}

		if _, err := conn.Exec(ctx, string(createSQL), t.FieldsVals()...); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to insert %T %v: %w", t, t, err).Error()
		}

		// TODO do we need to return it here?
		jsonBytes, err := json.Marshal(t)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to marshal the %T %v: %w", t, t, err).Error()
		}

		return http.StatusCreated, string(jsonBytes)
	})
}

type fieldsPtrsAware interface {
	FieldsPtrs() []any
}

func inst[T any]() T {
	return reflect.New(reflect.TypeFor[T]().Elem()).Interface().(T)
}

func Retrieve[T fieldsPtrsAware](ctx context.Context, retrieveSQL RetrieveSQL[T], id string) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		t := inst[T]()
		if err := conn.QueryRow(ctx, string(retrieveSQL), id).Scan(t.FieldsPtrs()...); err != nil {
			body := fmt.Errorf("failed to retrieve the %T with id %q not found: %w", t, id, err).Error()

			if errors.Is(err, pgx.ErrNoRows) {
				return http.StatusNotFound, body
			}

			if errors.Is(err, pgx.ErrTooManyRows) {
				return http.StatusConflict, body
			}

			return http.StatusInternalServerError, body
		}

		jsonBytes, err := json.Marshal(t)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to marshal the %T %v: %w", t, t, err).Error()
		}

		return http.StatusOK, string(jsonBytes)
	})
}

func Update[T fieldsValsAware](ctx context.Context, updateSQL UpdateSQL[T], t T) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		if err := t.Validate(); err != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to validate %v: %w", t, err).Error()
		}

		commandTag, err := conn.Exec(ctx, string(updateSQL), t.FieldsVals()...)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to update %T %v: %w", t, t, err).Error()
		}
		if commandTag.RowsAffected() == 0 {
			return http.StatusNotFound, fmt.Errorf("failed to update %T %v: %w", t, t, err).Error()
		}

		jsonBytes, err := json.Marshal(t)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to marshal the %T %v: %w", t, t, err).Error()
		}

		return http.StatusOK, string(jsonBytes)
	})
}

func Delete[T any](ctx context.Context, deleteSQL DeleteSQL[T], id string) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		commandTag, err := conn.Exec(ctx, string(deleteSQL), id)
		if err != nil {
			var t T // only to build the error
			return http.StatusInternalServerError, fmt.Errorf("failed to delete %T with id %v: %w", t, id, err).Error()
		}
		if commandTag.RowsAffected() == 0 {
			var t T // only to build the error
			return http.StatusNotFound, fmt.Errorf("failed to delete %T with id %v: no rows affected", t, id).Error()
		}

		return http.StatusNoContent, ""
	})
}

func List[T fieldsPtrsAware](ctx context.Context, listSQL ListSQL[T], offset int, limit int) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		sql := string(listSQL)

		rows, err := conn.Query(ctx, sql, offset, limit)
		if err != nil {
			body := fmt.Errorf("failed to execute list query for %T: %w", new(T), err).Error()
			return http.StatusInternalServerError, body
		}
		defer rows.Close()

		var buf bytes.Buffer

		buf.WriteByte('[')

		t := inst[T]()
		for rows.Next() {
			if err := rows.Scan(t.FieldsPtrs()...); err != nil {
				body := fmt.Errorf("failed to scan row for %T: %w", t, err).Error()
				return http.StatusInternalServerError, body
			}

			if buf.Len() > 0 {
				buf.WriteByte(',')
			}

			jsonBytes, err := json.Marshal(t)
			if err != nil {
				return http.StatusInternalServerError, fmt.Errorf("failed to marshal entity: %w", err).Error()
			}
			buf.Write(jsonBytes)
		}

		if err := rows.Err(); err != nil {
			body := fmt.Errorf("failed to iterate rows for %T: %w", t, err).Error()
			return http.StatusInternalServerError, body
		}

		buf.WriteByte(']')

		return http.StatusOK, buf.String()
	})
}
