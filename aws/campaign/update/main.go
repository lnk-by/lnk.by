package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/campaign"
)

func updateCampaign(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.Update(ctx, request, campaign.UpdateSQL, campaign.IdParam), nil
}

func main() {
	lambda.Start(updateCampaign)
}
