package api

import (
	"User_Service/internal/api/response"
	"User_Service/internal/erorrs"
	"User_Service/internal/service/model"
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// UpdateUser обновляет данные пользователя.
// @Summary Обновить данные пользователя
// @Description Обновляет данные пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body model.UpdateInfoUser true "Новые данные пользователя"
// @Success 200 {object} map[string]interface{} "Успешное обновление"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 401 {object} response.ErrorResponse "Неавторизованный доступ"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (r *Api) UpdateUser(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		r.logger.Error("Failed to get user", zap.Error(err))
		response.NewErrorResponse(c, http.StatusUnauthorized, "internal error")
		return
	}

	var input model.UpdateInfoUser

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = c.BindJSON(&input)
	if err != nil {
		r.logger.Error("Failed to binding JSON products", zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "bad request")
		return
	}

	id, err := r.service.User.UpdateUser(ctx, userId, input)
	if err != nil {
		r.logger.Error("Failed to update user", zap.Error(err))
		if errors.Is(err, erorrs.ErrUserNotFound) {
			r.logger.Error("User not found", zap.Error(err))
			response.NewErrorResponse(c, http.StatusBadRequest, "user not found")

		}
		response.NewErrorResponse(c, http.StatusBadRequest, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "successfuly updated", map[string]interface{}{
		"id": id,
	})
}
