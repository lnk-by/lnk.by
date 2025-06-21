package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/customer"
)

func createCustomer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	status, body := service.Create(ctx, customer.CreateSQL, []byte(request.Body))
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}, nil
}

func main() {
	lambda.Start(createCustomer)
}
