package user

import (
	"Gateway/internal/models/api"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type userSvc interface {
	RegisterUser(ctx context.Context, userReq api.CreateUserRequest) error
}

func RegisterUser(svc userSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx     = c.Request.Context()
			request api.CreateUserRequest
		)

		if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		if err := svc.RegisterUser(ctx, request); err != nil {
			log.WithField("login", request.Login).Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusCreated)
	}
}
