package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/organization"
)

func listOrganizations(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.List(ctx, request, organization.ListSQL), nil
}

func main() {
	lambda.Start(listOrganizations)
}
