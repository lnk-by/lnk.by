package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/auth"
	"github.com/lnk.by/shared/db"
	"net/http"
	"os"
)

func init() {
	if err := db.Init(context.Background(), os.Getenv("DB_URL"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	claims, err := auth.ParseJWT(request.Headers["Authorization"])
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	result, err := createShortURL(claims.OrganizationID, request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	body, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

func createShortURL(orgID string, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}
