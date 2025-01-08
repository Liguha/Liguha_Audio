package models

import (
	"context"
	log "github.com/sirupsen/logrus"
)

const (
	UserInfo = "userInfo"
)

type UserAuth struct {
	ID   uint64
	Role string
}

func GetUserFromContext(ctx context.Context) UserAuth {
	val := ctx.Value(UserInfo)
	user, ok := val.(UserAuth)
	if !ok {
		log.Fatalf("con not get users from context")
	}
	return user
}
