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

// SignUp регистрирует нового пользователя.
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.Login true "Данные для регистрации"
// @Success 200 {object} map[string]interface{} "Успешная регистрация"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 409 {object} response.ErrorResponse "Неверный email"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/sign-up [post]
func (r *Api) SignUp(c *gin.Context) {
	var input model.Login

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := c.BindJSON(&input)
	if err != nil {
		r.logger.Error("Failed to binding JSON", zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "bad request")
		return
	}

	id, err := r.service.Auth.CreateUser(ctx, input)
	if err != nil {
		r.logger.Error("Failed to create user", zap.Error(err))

		if errors.Is(err, erorrs.ErrUniqueConstraint) {
			r.logger.Error("User already exist", zap.Error(err))
			response.NewErrorResponse(c, http.StatusConflict, "email is already exist")
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "User created successfully", map[string]interface{}{
		"id": id,
	})

}

// SignIn аутентифицирует пользователя.
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает JWT-токен
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.Login true "Данные для входа"
// @Success 200 {object} map[string]interface{} "Успешный вход"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/sign-in [post]
func (r *Api) SignIn(c *gin.Context) {
	var input model.Login

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := c.BindJSON(&input)
	if err != nil {
		r.logger.Error("Failed to binding JSON", zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "bad request")
		return
	}

	token, err := r.service.Auth.GenerateJwtToken(ctx, input.Email, input.Password)
	if err != nil {
		r.logger.Error("Failed to generate token", zap.Error(err))

		if errors.Is(err, erorrs.ErrUserNotFound) {
			r.logger.Error("User not found", zap.Error(err))
			response.NewErrorResponse(c, http.StatusInternalServerError, "user not found")

		}
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "success", map[string]interface{}{
		"token": token,
	})
}
