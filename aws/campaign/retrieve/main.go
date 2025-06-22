package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/campaign"
)

func retrieveCampaign(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return adapter.Retrieve(ctx, request, campaign.RetrieveSQL, campaign.IdParam)
}

func main() {
	lambda.Start(retrieveCampaign)
}
