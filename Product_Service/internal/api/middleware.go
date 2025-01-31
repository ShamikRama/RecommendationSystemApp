package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Auth-Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing X-Auth-Token"})
			c.Abort()
			return
		}

		if token != "lwehvowhvowvhwovwfwefwefwefw" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: invalid X-Auth-Token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
