package response

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	slog.Error("Error occurred", "message", message)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{Message: message})
}
