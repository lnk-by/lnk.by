package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/customer"
	"github.com/lnk.by/shared/utils"
)

func createCustomerFromUser(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	id := event.Request.UserAttributes["sub"]
	email := event.Request.UserAttributes["email"]

	c := customer.Customer{
		ID:     id,
		Email:  email,
		Name:   event.UserName,
		Status: utils.StatusActive,
	}

	status, message := service.CreateRecord(ctx, customer.CreateSQL, &c)

	if status != http.StatusCreated {
		return event, fmt.Errorf("failed to create customer for user id: %s, name: %s, email: %s, error: %s", id, event.UserName, email, message)
	}

	return event, nil

}

func main() {
	adapter.LambdaMain(createCustomerFromUser)
}
