package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/user"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var u user.User // TODO build user from request
	status, body := service.Create(ctx, user.Ops, &u)
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func main() {
	lambda.Start(handleRequest)
}
