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

type Ops[T any] struct {
	Create   string
	Retrieve string
	Update   string
	Delete   string
}

//func typeName(v any) entity {
//	t := fmt.Sprintf("%T", v)
//	for strings.HasPrefix(t, "*") {
//		t = t[1:]
//	}
//	return entity(t)
//}
//
//type entity string
//
//var (
//	entUser         entity = typeName(user.User{})
//	organization entity = "organization.Organization"
//	campaign     entity = "campaign.Campaign"
//)
//
//type operation string
//
//const (
//	opCreate   operation = "create"
//	opRetrieve operation = "retrieve"
//	opUpdate   operation = "update"
//	opDelete   operation = "delete"
//)

//var sqlStatements = map[entity]map[operation]string{}

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

func Create[T fieldsValsAware](ctx context.Context, ops Ops[T], t T) (int, string) {
	typeOfT := fmt.Sprintf("%T", t)
	fmt.Println(typeOfT)

	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		if _, err := conn.Exec(ctx, ops.Create, t.FieldsVals()...); err != nil {
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

func Retrieve[T fieldsPtrsAware](ctx context.Context, ops Ops[T], id string) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		var t T
		if err := conn.QueryRow(ctx, ops.Retrieve, id).Scan(t.FieldsPtrs()...); err != nil {
			body := fmt.Errorf("failed to retrieve the %T with id %q not found: %w", t, id, err).Error()

			if errors.Is(err, pgx.ErrNoRows) {
				return http.StatusNotFound, body
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

func Update[T fieldsValsAware](ctx context.Context, ops Ops[T], t T) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		commandTag, err := conn.Exec(ctx, ops.Update, t.FieldsVals()...)
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

func Delete[T any](ctx context.Context, ops Ops[T], id string) (int, string) {
	return withConn(ctx, func(conn *pgxpool.Conn) (int, string) {
		commandTag, err := conn.Exec(ctx, ops.Delete, id)
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
