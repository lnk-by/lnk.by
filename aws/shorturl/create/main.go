package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service"
	shorturlservice "github.com/lnk.by/shared/service/shorturl"
)

func createShortURL(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := service.GetUUIDFromAuthorization(request.Headers["authorization"])
	slog.Info("Creating short URL", "userID", userID)
	status, body := shorturlservice.CreateShortURL(ctx, []byte(request.Body), userID)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: adapter.StandardHeaders}, nil
}

func main() {
	adapter.LambdaMain(createShortURL)
}
