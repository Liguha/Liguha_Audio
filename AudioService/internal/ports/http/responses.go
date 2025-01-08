package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIResponse[T any] struct {
	Data *T `json:"data,omitempty"`
}

func SuccessResponse[T any](c *gin.Context, data *T) {
	response := APIResponse[T]{Data: data}

	c.JSON(http.StatusOK, response)
}
