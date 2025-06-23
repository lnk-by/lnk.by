package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/customer"
)

func retrieveCustomer(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	return adapter.Retrieve(ctx, request, customer.RetrieveSQL, customer.IdParam)
}

func main() {
	lambda.Start(retrieveCustomer)
}
