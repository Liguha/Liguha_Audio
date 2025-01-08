package middleware

import (
	"AudioService/internal/models"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	userIDHeader   = "x-user-id"
	userRoleHeader = "x-user-role"
)

func SimpleMiddleware(c *gin.Context) {
	userID := c.GetHeader(userIDHeader)
	userRole := c.GetHeader(userRoleHeader)

	userIDInt, _ := strconv.Atoi(userID)

	userInfo := models.UserAuth{
		ID:   uint64(userIDInt),
		Role: userRole,
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), models.UserInfo, userInfo))
	c.Next()
}
