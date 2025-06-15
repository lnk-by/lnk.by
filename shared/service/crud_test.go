package service

import (
	"context"
	"fmt"
	"github.com/lnk.by/shared/service/user"
	"testing"
)

func TestCreate(t *testing.T) {
	status, body := Create[*user.User](context.Background(), user.Ops, &user.User{})
	fmt.Println(status, body)
}
