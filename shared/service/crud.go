package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
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

func Parse[T Creatable](ctx context.Context, content []byte) (T, error) {
	t := inst[T]()
	if err := json.Unmarshal(content, t); err != nil {
		return t, fmt.Errorf("failed to unmarshal %T from JSON: %w", t, err)
	}
	if err := t.Validate(); err != nil {
		return t, fmt.Errorf("failed to validate %T: %w", t, err)
	}
	return t, nil
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

type Identifiable[K any] interface {
	ParseID(string) (K, error)
}

type Retrievable[K any] interface {
	Identifiable[K]
	FieldsPtrs() []any
}

func inst[T any]() T {
	return reflect.New(reflect.TypeFor[T]().Elem()).Interface().(T)
}

func Retrieve[K any, T Retrievable[K]](ctx context.Context, retrieveSQL RetrieveSQL[T], id string, transformer func(t T) (T, error)) (int, string) {
	return marshal(retrieve(ctx, retrieveSQL, id, transformer))
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

func RetrieveValueAndMarshalError[K any, T Retrievable[K]](ctx context.Context, retrieveSQL RetrieveSQL[T], idString string, args ...string) (int, T, string) {
	status, value, err := retrieve(ctx, retrieveSQL, idString, func(t T) (T, error) { return t, nil }, args...)
	if err != nil {
		_, strErr := failed(status, err)
		return status, value, strErr
	}

	return status, value, ""
}

func retrieve[K any, T Retrievable[K]](ctx context.Context, retrieveSQL RetrieveSQL[T], idString string, transformer func(t T) (T, error), args ...string) (int, T, error) {
	t := inst[T]()
	id, err := t.ParseID(idString)
	if err != nil {
		return http.StatusNotFound, t, fmt.Errorf("failed to parse %T ID: %v: %w", t, idString, err)
	}

	return withConn(ctx, func(conn *pgxpool.Conn) (int, T, error) {
		var sql string
		if len(args) > 0 {
			anyArgs := make([]any, len(args))
			for i, v := range args {
				anyArgs[i] = v
			}
			sql = fmt.Sprintf(string(retrieveSQL), anyArgs...)
		} else {
			sql = string(retrieveSQL)
		}
		if err := conn.QueryRow(ctx, sql, id).Scan(t.FieldsPtrs()...); err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				return http.StatusNotFound, t, fmt.Errorf("failed to retrieve the %T with id '%v': %w", t, id, err)
			case errors.Is(err, pgx.ErrTooManyRows):
				return http.StatusConflict, t, fmt.Errorf("failed to retrieve the %T with id '%v': %w", t, id, err)
			default:
				return http.StatusInternalServerError, t, fmt.Errorf("failed to retrieve the %T with id '%v': %w", t, id, err)
			}
		}
		t, err = transformer(t)
		if err != nil {
			return http.StatusInternalServerError, t, fmt.Errorf("failed to transform row value %T: %w", t, err)
		}
		return http.StatusOK, t, nil
	})
}

type Updatable[K any] interface {
	Identifiable[K]
	FieldsValsAware
	WithID(K)
}

func UpdateFromReqBody[K any, T Updatable[K]](ctx context.Context, updateSQL UpdateSQL[T], idString string, body io.ReadCloser, finalizer func(id K, t T) error) (int, string) {
	content, err := io.ReadAll(body)
	if err != nil {
		return failed(http.StatusInternalServerError, fmt.Errorf("failed to read request body: %w", err))
	}
	return Update(ctx, updateSQL, idString, content, finalizer)
}

func Update[K any, T Updatable[K]](ctx context.Context, updateSQL UpdateSQL[T], idString string, content []byte, finalizer func(id K, t T) error) (int, string) {
	t := inst[T]()
	if err := json.Unmarshal(content, t); err != nil {
		return failed(http.StatusBadRequest, fmt.Errorf("failed to unmarshal %T from JSON: %w", t, err))
	}
	if err := t.Validate(); err != nil {
		return failed(http.StatusBadRequest, fmt.Errorf("failed to validate %T: %w", t, err))
	}

	id, err := t.ParseID(idString)
	if err != nil {
		return failed(http.StatusNotFound, fmt.Errorf("failed to parse %T ID: %v: %w", t, idString, err))
	}

	t.WithID(id)

	return marshal(withConn(ctx, func(conn *pgxpool.Conn) (int, T, error) {
		commandTag, err := conn.Exec(ctx, string(updateSQL), t.FieldsVals()...)
		switch {
		case err != nil:
			return http.StatusInternalServerError, t, fmt.Errorf("failed to update %T %v: %w", t, t, err)
		case commandTag.RowsAffected() == 0:
			return http.StatusNotFound, t, fmt.Errorf("failed to update %T %v: %w", t, t, pgx.ErrNoRows)
		case commandTag.RowsAffected() > 1: // TODO is it too late?
			return http.StatusConflict, t, fmt.Errorf("failed to update %T %v: %w", t, t, pgx.ErrTooManyRows)
		}

		if err = finalizer(id, t); err != nil {
			return http.StatusInternalServerError, t, fmt.Errorf("failed to finilize updating of %T with id %v: %w", t, id, err)
		}

		return http.StatusOK, t, nil
	}))
}

func Delete[K any, T Identifiable[K]](ctx context.Context, deleteSQL DeleteSQL[T], idString string, finalizer func(id K) error) (int, string) {
	var t T
	id, err := t.ParseID(idString) // it it OK if t is nil
	if err != nil {
		return failed(http.StatusNotFound, fmt.Errorf("failed to parse %T ID: %v: %w", t, idString, err))
	}

	return marshal(withConn(ctx, func(conn *pgxpool.Conn) (int, T, error) {
		commandTag, err := conn.Exec(ctx, string(deleteSQL), id)
		switch {
		case err != nil:
			return http.StatusInternalServerError, t, fmt.Errorf("failed to delete %T with id %v: %w", t, id, err)
		case commandTag.RowsAffected() == 0:
			return http.StatusNotFound, t, fmt.Errorf("failed to delete %T with id %v: %w", t, id, pgx.ErrNoRows)
		case commandTag.RowsAffected() > 1: // TODO is it too late?
			return http.StatusNotFound, t, fmt.Errorf("failed to delete %T with id %v: %w", t, id, pgx.ErrTooManyRows)
		}

		if err = finalizer(id); err != nil {
			return http.StatusInternalServerError, t, fmt.Errorf("failed to finilize deletion of %T with id %v: %w", t, id, err)
		}

		return http.StatusNoContent, t, nil
	}))
}

func List[K any, T Retrievable[K]](ctx context.Context, listSQL ListSQL[T], userID *uuid.UUID, offset int, limit int, transformer func(t T) (T, error)) (int, string) {
	return marshal(withConn(ctx, func(conn *pgxpool.Conn) (int, []T, error) {
		sql := string(listSQL)
		slog.Info("List 1", "sql", sql)

		rows, err := conn.Query(ctx, sql, userID, offset, limit)
		slog.Info("List 2")
		if err != nil {
			slog.Info("List 2.1", "error", err)
			return http.StatusInternalServerError, nil, fmt.Errorf("failed to execute list query for %T: %w", new(T), err)
		}
		defer rows.Close()
		slog.Info("List 3")

		ts := make([]T, 0)
		for rows.Next() {
			t := inst[T]()
			slog.Info("List 4")
			if err := rows.Scan(t.FieldsPtrs()...); err != nil {
				slog.Info("List 4.1", "error", err)
				return http.StatusInternalServerError, ts, fmt.Errorf("failed to scan row to build the %T: %w", t, err)
			}
			slog.Info("List 5")
			t, err = transformer(t)
			slog.Info("List 6")
			if err != nil {
				slog.Info("List 6.1", "error", err)
				return http.StatusInternalServerError, ts, fmt.Errorf("failed to transform row value %T: %w", t, err)
			}
			slog.Info("List 7")
			ts = append(ts, t)
			slog.Info("List 8")
		}
		slog.Info("List 9")

		return http.StatusOK, ts, nil
	}))
}

func UUID() uuid.UUID {
	return uuid.Must(uuid.NewV1())
}

func ToUUID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}
	id, err := uuid.FromString(s)
	if err != nil {
		return nil
	}
	return &id
}

func GetUUIDFromAuthorization(authHeader string) *uuid.UUID {
	claims := getClaimsFromAuthorization(authHeader)
	if claims != nil {
		if subStr, ok := claims["sub"].(string); ok {
			return ToUUID(subStr)
		}
	}
	return nil
}

func getClaimsFromAuthorization(authHeader string) map[string]interface{} {
	var claims map[string]interface{}

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		parts := strings.Split(token, ".")
		if len(parts) == 3 {
			payloadSegment := parts[1]

			// Pad to make it valid base64
			missingPadding := len(payloadSegment) % 4
			if missingPadding != 0 {
				payloadSegment += strings.Repeat("=", 4-missingPadding)
			}

			decoded, err := base64.RawURLEncoding.DecodeString(payloadSegment)
			if err != nil {
				slog.Error("Failed to decode JWT payload:", "error", err)
			} else {
				err = json.Unmarshal(decoded, &claims)
				if err != nil {
					fmt.Println("Failed to parse JWT claims:", err)
				}
			}
		} else {
			slog.Error("Invalid JWT format")
		}
	}
	return claims
}

var (
	ErrNameRequired      = errors.New("name is required")
	ErrIDManagedByServer = errors.New("ID is managed by the server")
)
