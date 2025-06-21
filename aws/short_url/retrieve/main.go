package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/short_url"
)

func retrieveShortURL(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	status, body := service.Retrieve(ctx, short_url.RetrieveSQL, request.PathParameters[short_url.IdParam])
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func main() {
	lambda.Start(retrieveShortURL)
}
