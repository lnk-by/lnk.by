package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/shared/service"
)

func Create[T service.Creatable](ctx context.Context, request events.APIGatewayProxyRequest, sql service.CreateSQL[T]) events.APIGatewayProxyResponse {
	status, body := service.Create(ctx, sql, []byte(request.Body))
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func Retrieve[T service.FieldsPtrsAware](ctx context.Context, request events.APIGatewayProxyRequest, sql service.RetrieveSQL[T], idParam string) events.APIGatewayProxyResponse {
	status, body := service.Retrieve(ctx, sql, request.PathParameters[idParam])
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func Update[T service.Updatable](ctx context.Context, request events.APIGatewayProxyRequest, sql service.UpdateSQL[T], idParam string) events.APIGatewayProxyResponse {
	status, body := service.Update(ctx, sql, request.PathParameters[idParam], []byte(request.Body))
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func Delete[T service.FieldsValsAware](ctx context.Context, request events.APIGatewayProxyRequest, sql service.DeleteSQL[T], idParam string) events.APIGatewayProxyResponse {
	status, body := service.Delete(ctx, sql, request.PathParameters[idParam])
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func List[T service.FieldsPtrsAware](ctx context.Context, request events.APIGatewayProxyRequest, listSQL service.ListSQL[T]) events.APIGatewayProxyResponse {
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
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
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

func badRequestResponse(err error) events.APIGatewayProxyResponse {
	slog.Warn("BadRequest:", "error", err)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}
}
