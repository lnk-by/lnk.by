package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/organization"
)

func retrieveOrganization(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return adapter.Retrieve(ctx, request, organization.RetrieveSQL, organization.IdParam)
}

func main() {
	lambda.Start(retrieveOrganization)
}
