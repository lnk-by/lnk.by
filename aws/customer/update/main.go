package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/customer"
)

func updateCustomer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.Update(ctx, request, customer.UpdateSQL, customer.IdParam), nil

}

func main() {
	lambda.Start(updateCustomer)
}
