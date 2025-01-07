package user

import (
	"Gateway/internal/entity"
	httpHandler "Gateway/internal/ports/http"
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"Gateway/internal/models/api"
)

type userLoginSvc interface {
	Login(ctx context.Context, userReq api.CreateUserRequest) (string, error)
}

func LoginUser(svc userLoginSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx  = c.Request.Context()
			auth = c.GetHeader("authorization")
		)

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ErrUnauthorized)
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ErrUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ErrUnauthorized)
			return
		}

		decodedParts := strings.Split(string(decoded), ":")
		userReq := api.CreateUserRequest{
			Login:    decodedParts[0],
			Password: decodedParts[1],
		}

		token, err := svc.Login(ctx, userReq)
		if err != nil {
			log.WithField("login", userReq.Login).Error(err)
			switch {
			case errors.Is(err, entity.ErrNotFound):
				c.AbortWithStatusJSON(http.StatusNotFound, err)
			case errors.Is(err, entity.ErrUnauthorized):
				c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			default:
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		httpHandler.SuccessResponse(c, &token)
	}
}
