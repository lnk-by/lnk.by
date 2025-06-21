package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/short_url"
)

func updateShortURL(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	status, body := service.Update(ctx, short_url.UpdateSQL, request.PathParameters[short_url.IdParam], []byte(request.Body))
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}, nil
}

func main() {
	lambda.Start(updateShortURL)
}
