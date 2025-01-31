package api

import (
	"Analytics_Service/internal/api/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserActions возвращает список действий пользователя.
// @Summary Получить действия пользователя
// @Description Возвращает список всех действий пользователя
// @Tags user
// @Produce json
// @Success 200 {object} map[string]interface{} "Успешный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /stats/actions [get]
func (h *Api) GetUserActions(c *gin.Context) {
	actions, err := h.service.GetUserActions()
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "success", map[string]interface{}{
		"actions": actions,
	})
}

// GetUserEvents возвращает список событий пользователя.
// @Summary Получить события пользователя
// @Description Возвращает список всех событий пользователя
// @Tags user
// @Produce json
// @Success 200 {object} map[string]interface{} "Успешный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /stats/events [get]
func (h *Api) GetUserEvents(c *gin.Context) {
	events, err := h.service.GetUserEvents()
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "success", map[string]interface{}{
		"actions": events,
	})
}

// GetStats возвращает статистику продуктов.
// @Summary Получить статистику продуктов
// @Description Возвращает статистику по продуктам
// @Tags stats
// @Produce json
// @Success 200 {object} map[string]interface{} "Успешный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /stats/products [get]
func (h *Api) GetStats(c *gin.Context) {
	stats, err := h.service.GetProductStats()
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "success", map[string]interface{}{
		"actions": stats,
	})
}
