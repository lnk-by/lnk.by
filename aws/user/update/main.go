package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/user"
)

func UpdateUser(ctx context.Context, userId string, requestedUser user.User) (int, string) {
	if requestedUser.Email == "" || requestedUser.Name == "" {
		return http.StatusBadRequest, "Email and Name are required"
	}
	if requestedUser.ID != "" {
		return http.StatusBadRequest, "User ID is managed by the server"
	}
	return service.Update(ctx, user.UpdateSQL, &requestedUser)
}

func main() {
	lambda.Start(common.UpdateAdapter(UpdateUser, user.UserIdParam, common.StringIDParser))
}
