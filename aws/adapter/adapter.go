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

var standardHeaders = map[string]string{
	"Content-Type":                "application/json",
	"Access-Control-Allow-Origin": "*",
}

func Create[T service.Creatable](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.CreateSQL[T]) events.APIGatewayV2HTTPResponse {
	status, body := service.Create(ctx, sql, []byte(request.Body))
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func Retrieve[T service.FieldsPtrsAware](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.RetrieveSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	status, body := service.Retrieve(ctx, sql, request.PathParameters[idParam])
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func Update[T service.Updatable](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.UpdateSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	status, body := service.Update(ctx, sql, request.PathParameters[idParam], []byte(request.Body))
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func Delete[T service.FieldsValsAware](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.DeleteSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	status, body := service.Delete(ctx, sql, request.PathParameters[idParam])
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func List[T service.FieldsPtrsAware](ctx context.Context, request events.APIGatewayV2HTTPRequest, listSQL service.ListSQL[T]) events.APIGatewayV2HTTPResponse {
	params := request.QueryStringParameters
	offset, err := parseQueryInt(params, "offset", 0)
	if err != nil {
		return badRequestResponse(err)
	}
	limit, err := parseQueryInt(params, "limit", math.MaxInt32)
	if err != nil {
		return badRequestResponse(err)
	}
	status, body := service.List(ctx, listSQL, offset, limit)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
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
	if err := db.InitFromEnvironement(ctx); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}
