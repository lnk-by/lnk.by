package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/short_url"
)

/*
attempts := 10
	if requestedShortURL.Key != "" {
		attempts = 1
		key := requestedShortURL.Key
		short_url.Custom.Put(key, true)
		defer short_url.Custom.Delete(key)
	}
*/

func createShortURL(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	status, body := service.Create(ctx, short_url.CreateSQL, []byte(request.Body))
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}, nil
}

func main() {
	lambda.Start(createShortURL)
}
