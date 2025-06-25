package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/campaign"
)

func retrieveCampaign(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	return adapter.Retrieve(ctx, request, campaign.RetrieveSQL, campaign.IdParam)
}

func main() {
	adapter.LambdaMain(retrieveCampaign)
}
