package api

import (
	"Recommendation_Service/internal/api/response"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// GetRecommend возвращает рекомендации для пользователя.
// @Summary Получить рекомендации
// @Description Возвращает список рекомендаций для пользователя по его ID
// @Tags recommendations
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]interface{} "Успешный запрос"
// @Failure 400 {object} response.ErrorResponse "Неверный ID пользователя"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /rec/{id} [get]
func (r *RecApi) GetRecommend(c *gin.Context) {
	id := c.Param("id")

	userId, err := strconv.Atoi(id)
	if err != nil {
		r.logger.Error("Failed to get id",
			zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	recommendations, err := r.service.GetRecommendations(userId)
	if err != nil {
		r.logger.Error("Failed to get recommendation",
			zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "success", map[string]interface{}{
		"recommendations": recommendations,
	})
}
