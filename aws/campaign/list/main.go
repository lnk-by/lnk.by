package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service/customer"
)

func listCustomers(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return common.ListEntities(ctx, request, customer.ListSQL)
}

func main() {
	lambda.Start(listCustomers)
}
