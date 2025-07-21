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
	"github.com/lnk.by/shared/service/shorturl"
	"github.com/lnk.by/shared/service/stats"
)

var standardHeaders = map[string]string{
	"Content-Type":                "application/json",
	"Access-Control-Allow-Origin": "*",
}

func CreateShortURL(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	url, err := service.Parse[*shorturl.ShortURL](ctx, []byte(request.Body))
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusBadRequest, Body: http.StatusText(http.StatusBadRequest), Headers: standardHeaders}
	}

	e := stats.Event{Key: url.Key}
	createSQLs := []service.CreateSQL[*stats.Event]{stats.CreateTotalSQL, stats.CreateDailySQL, stats.CreateHourlySQL, stats.CreateUserAgentSQL, stats.CreateCountrySQL}
	for _, sql := range createSQLs {
		status, body := service.CreateRecord(ctx, sql, &e, 0)
		if status >= http.StatusBadRequest {
			return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
		}
	}

	status, body := service.CreateRecord(ctx, shorturl.CreateSQL, url, 0)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func Create[T service.Creatable](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.CreateSQL[T]) events.APIGatewayV2HTTPResponse {
	status, body := service.Create(ctx, sql, []byte(request.Body))
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func Retrieve[K any, T service.Retrievable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.RetrieveSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	status, body := service.Retrieve(ctx, sql, request.PathParameters[idParam])
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func Update[K any, T service.Updatable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.UpdateSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	status, body := service.Update(ctx, sql, request.PathParameters[idParam], []byte(request.Body))
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func Delete[K any, T service.Identifiable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, sql service.DeleteSQL[T], idParam string) events.APIGatewayV2HTTPResponse {
	status, body := service.Delete(ctx, sql, request.PathParameters[idParam])
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}
}

func List[K any, T service.Retrievable[K]](ctx context.Context, request events.APIGatewayV2HTTPRequest, listSQL service.ListSQL[T]) events.APIGatewayV2HTTPResponse {
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
	if err := db.InitFromEnvironment(ctx); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}
