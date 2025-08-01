package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/aws/s3client"
	"github.com/lnk.by/shared/service/landingpage"
)

func retrieveLandingPageTemplate(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	template, err := landingpage.RetrieveLandingPageTemplate(ctx, request.PathParameters["name"])
	bytes, err := json.Marshal(template)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: err.Error(), Headers: adapter.StandardHeaders}, nil
	}
	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusOK, Body: string(bytes), Headers: adapter.StandardHeaders}, nil
}

func main() {
	s3client.InitializeFromEnvironment(context.Background())
	adapter.LambdaMain(retrieveLandingPageTemplate)
}
