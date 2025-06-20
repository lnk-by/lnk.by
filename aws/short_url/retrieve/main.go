package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/short_url"
)

func retrieveShortURL(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return adapter.Retrieve(ctx, request, short_url.RetrieveSQL, short_url.IdParam)
}

func main() {
	lambda.Start(retrieveShortURL)
}
