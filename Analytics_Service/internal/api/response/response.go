package response

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func NewResponse(c *gin.Context, statusCode int, message string, data map[string]interface{}) {
	slog.Info("Success", "message", message)
	c.JSON(statusCode, Response{Message: message, Data: data})
}
