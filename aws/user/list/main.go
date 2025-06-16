package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/user"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	params := request.QueryStringParameters
	limit, err := parseQueryInt(params, "limit", 0)
	if err != nil {
		return badRequestResponse(err)
	}
	offset, err := parseQueryInt(params, "offset", 0)
	if err != nil {
		return badRequestResponse(err)
	}
	status, body := service.List(ctx, user.ListSQL, limit, offset)
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func badRequestResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}
}

func parseQueryInt(params map[string]string, key string, defaultValue int) (int, error) {
	valStr, ok := params[key]
	if !ok || valStr == "" {
		return defaultValue, nil
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value for '%s': '%s' is not an integer", key, valStr)
	}
	return val, nil
}

func main() {
	lambda.Start(handleRequest)
}
