package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/organization"
)

func deleteOrganization(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.Delete(ctx, request, organization.DeleteSQL, organization.IdParam), nil
}

func main() {
	lambda.Start(deleteOrganization)
}
