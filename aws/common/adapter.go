package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log/slog"
)

var logger = slog.Default()

type CustomCreateHandler[T any] func(ctx context.Context, input T) (int, string)
type CustomUpdateHandler[T any, ID any] func(ctx context.Context, id ID, input T) (int, string)

func CreateAdapter[T any](handler CustomCreateHandler[T]) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var input T
		err := json.Unmarshal([]byte(request.Body), &input)
		if err != nil {
			logger.Warn("Failed to parse request", "error", err)
			return badRequestResponse(err), nil
		}

		status, body := handler(ctx, input)

		return events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       body,
		}, nil
	}
}

func UpdateAdapter[T any, ID any](
	handler func(ctx context.Context, id ID, input T) (int, string),
	pathParamName string,
	idParser func(string) (ID, error),
) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var input T
		err := json.Unmarshal([]byte(request.Body), &input)
		if err != nil {
			return badRequestResponse(fmt.Errorf("failed to parse request: %w", err)), nil
		}

		idStr, ok := request.PathParameters[pathParamName]
		if !ok {
			return badRequestResponse(fmt.Errorf("parameter %q does not exist", pathParamName)), nil
		}
		id, err := idParser(idStr)
		if err != nil {
			return badRequestResponse(err), nil
		}

		status, body := handler(ctx, id, input)

		return events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       body,
		}, nil
	}
}

func StringIDParser(s string) (string, error) {
	return s, nil
}
