package api

import (
	"User_Service/internal/api/response"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"
	userCtx    = "user_id"
)

func (r *Api) UserIdentity(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		r.logger.Info("Empty header")
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		r.logger.Info("Wrong header")
		response.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	token := headerParts[1]
	if token == "" {
		r.logger.Info("Empty token")
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty token")
		return
	}

	userId, err := r.service.Auth.ParseToken(token)
	if err != nil {
		r.logger.Error("Failed to parse the token")
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id not int")
	}

	return idInt, nil

}
