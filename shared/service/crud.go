package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lnk.by/shared/db"
	"net/http"
)

type CreateSQL[T any] string
type RetrieveSQL[T any] string
type UpdateSQL[T any] string
type DeleteSQL[T any] string

func withConn(ctx context.Context, f func(conn *pgxpool.Conn) (int, string)) (int, string) {
	conn, err := db.Get(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get DB connection: %w", err).Error()
	}
	defer conn.Release()

	return f(conn)
}

type fieldsValsAware interface {
	FieldsVals() []any
}

func Create[T fieldsValsAware](ctx context.Context, createSQL CreateSQL[T], t T) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		if _, err := conn.Exec(ctx, string(createSQL), t.FieldsVals()...); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to insert %T %v: %w", t, t, err).Error()
		}

		// TODO do we need to return it here?
		bytes, err := json.Marshal(t)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to marshal the %T %v: %w", t, t, err).Error()
		}

		return http.StatusCreated, string(bytes)
	})
}

type fieldsPtrsAware interface {
	FieldsPtrs() []any
}

func Retrieve[T fieldsPtrsAware](ctx context.Context, retrieveSQL RetrieveSQL[T], id string) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		var t T
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

		bytes, err := json.Marshal(t)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to marshal the %T %v: %w", t, t, err).Error()
		}

		return http.StatusOK, string(bytes)
	})
}

func Update[T fieldsValsAware](ctx context.Context, updateSQL UpdateSQL[T], t T) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		commandTag, err := conn.Exec(ctx, string(updateSQL), t.FieldsVals()...)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to update %T %v: %w", t, t, err).Error()
		}
		if commandTag.RowsAffected() == 0 {
			return http.StatusNotFound, fmt.Errorf("failed to update %T %v: %w", t, t, err).Error()
		}

		// TODO do we need to return it here?
		bytes, err := json.Marshal(t)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to marshal the %T %v: %w", t, t, err).Error()
		}

		return http.StatusOK, string(bytes)
	})
}

func Delete[T any](ctx context.Context, deleteSQL DeleteSQL[T], id string) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		commandTag, err := conn.Exec(ctx, string(deleteSQL), id)
		if err != nil {
			var t T // only to build the error
			return http.StatusInternalServerError, fmt.Errorf("failed to delete %T with id %v: %w", t, t, err).Error()
		}
		if commandTag.RowsAffected() == 0 {
			var t T // only to build the error
			return http.StatusNotFound, fmt.Errorf("failed to delete %T %v: %w", t, t, err).Error()
		}

		return http.StatusNoContent, ""
	})
}
