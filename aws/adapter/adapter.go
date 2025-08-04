package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/db"
	"github.com/lnk.by/shared/service"
)

var StandardHeaders = map[string]string{
	"Content-Type":                "application/json",
	"Access-Control-Allow-Origin": "*",
}

func Create[T service.Creatable](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.CreateSQL[T]) events.APIGatewayV2HTTPResponse {
	status, body := service.Create(ctx, sql, []byte(request.Body))
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: StandardHeaders}
}

func Retrieve[K any, T service.Retrievable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.RetrieveSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	return RetrieveAndTransform(ctx, request, sql, idParam, func(t T) (T, error) { return t, nil })
}

func RetrieveAndTransform[K any, T service.Retrievable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.RetrieveSQL[T], idParam string, transformer func(t T) (T, error)) events.APIGatewayV2HTTPResponse {
	status, body := service.Retrieve(ctx, sql, request.PathParameters[idParam], transformer)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: StandardHeaders}
}

func Update[K any, T service.Updatable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.UpdateSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	return UpdateAndFinalize(ctx, request, sql, idParam, func(id K, t T) error { return nil })
}

func UpdateAndFinalize[K any, T service.Updatable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.UpdateSQL[T], idParam string, finalizer func(id K, t T) error) events.APIGatewayV2HTTPResponse {
	status, body := service.Update(ctx, sql, request.PathParameters[idParam], []byte(request.Body), finalizer)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: StandardHeaders}
}

func Delete[K any, T service.Identifiable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.DeleteSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	return DeleteAndFinalize(ctx, request, sql, idParam, func(id K) error { return nil })
}

func DeleteAndFinalize[K any, T service.Identifiable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.DeleteSQL[T], idParam string, finalizer func(id K) error) events.APIGatewayV2HTTPResponse {
	status, body := service.Delete(ctx, sql, request.PathParameters[idParam], finalizer)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: StandardHeaders}
}

func List[K any, T service.Retrievable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, listSQL service.ListSQL[T]) events.APIGatewayV2HTTPResponse {
	return ListAndTransform(ctx, request, listSQL, func(t T) (T, error) { return t, nil })
}

func ListAndTransform[K any, T service.Retrievable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, listSQL service.ListSQL[T], transformer func(t T) (T, error)) events.APIGatewayV2HTTPResponse {
	params := request.QueryStringParameters
	offset, err := parseQueryInt(params, "offset", 0)
	if err != nil {
		return badRequestResponse(err)
	}
	limit, err := parseQueryInt(params, "limit", math.MaxInt32)
	if err != nil {
		return badRequestResponse(err)
	}
	userID := service.ToUUID(request.RequestContext.Authorizer.JWT.Claims["sub"])
	status, body := service.List(ctx, listSQL, userID, offset, limit, transformer)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: StandardHeaders}
}

func parseQueryInt(params map[string]string, key string, defaultValue int) (int, error) {
	valStr, ok := params[key]
	if !ok || valStr == "" {
		return defaultValue, nil
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("failed to invalid value for %q: %q is not an integer", key, valStr)
	}
	return val, nil
}

func badRequestResponse(err error) events.APIGatewayV2HTTPResponse {
	slog.Warn("BadRequest:", "error", err)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}
}

func LambdaMain(handler interface{}) {
	ctx := context.Background()
	if err := db.InitFromEnvironment(ctx); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}
