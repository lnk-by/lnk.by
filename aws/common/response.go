package common

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func badRequestResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}
}
