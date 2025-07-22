package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service"
	shorturlservice "github.com/lnk.by/shared/service/shorturl"
)

var standardHeaders = map[string]string{
	"Content-Type":                "application/json",
	"Access-Control-Allow-Origin": "*",
}

func createShortURL(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := service.GetUUIDFromAuthorization(request.Headers["authorization"])
	status, body := shorturlservice.CreateShortURL(ctx, []byte(request.Body), userID)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: standardHeaders}, nil
}

func main() {
	adapter.LambdaMain(createShortURL)
}
