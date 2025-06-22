package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/customer"
)

func deleteCustomer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.Delete(ctx, request, customer.DeleteSQL, customer.IdParam), nil
}

func main() {
	lambda.Start(deleteCustomer)
}
