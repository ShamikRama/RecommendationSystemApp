package api

import (
	"User_Service/internal/api/response"
	"User_Service/pkg/product"
	"context"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// GetProductsHandler возвращает список продуктов.
// @Summary Получить список продуктов
// @Description Возвращает список всех продуктов с поддержкой пагинации
// @Tags product
// @Produce json
// @Param page query int false "Номер страницы (начинается с 0). По умолчанию 0, если не указан."
// @Param page_size query int false "Количество продуктов на странице. По умолчанию 0 (все продукты), если не указан."
// @Success 200 {object} map[string]interface{} "Успешный запрос. Возвращает список продуктов и метаданные пагинации."
// @Failure 400 {object} response.ErrorResponse "Неверные параметры запроса (например, page или page_size не являются числами)."
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера."
// @Router /products/ [get]
func (r *Api) GetProductsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Извлекаем параметры пагинации
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "0"))
	if err != nil {
		r.logger.Error("Invalid page_size parameter",
			zap.Error(err),
			zap.String("page_size", c.Query("page_size")))
		response.NewErrorResponse(c, http.StatusBadRequest, "bad request")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		r.logger.Error("Invalid page parameter",
			zap.Error(err),
			zap.String("page", c.Query("page")))
		response.NewErrorResponse(c, http.StatusBadRequest, "bad request")
		return
	}

	products, err := r.product.GetProducts(ctx, pageSize, page)
	if err != nil {
		r.logger.Error("Failed to get products",
			zap.Error(err),
			zap.Int("page_size", pageSize),
			zap.Int("page", page))
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, products)
}

// CreateCartHandler создает корзину.
// @Summary Создать корзину
// @Description Создает корзину для пользователя
// @Tags product
// @Accept json
// @Produce json
// @Param input body product.CartCreateRequest true "Данные для создания корзины"
// @Success 200 {object} map[string]interface{} "Успешное создание"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/ [post]
func (r *Api) CreateCartHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var request product.CartCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.logger.Error("Failed to binding JSON", zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "bad request")
		return
	}

	err := r.product.CreateCart(ctx, request)
	if err != nil {
		r.logger.Error("Failed to create the cart", zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "cart created successfully", nil)
}

// DeleteCartHandler удаляет корзину.
// @Summary Удалить корзину
// @Description Удаляет корзину пользователя
// @Tags product
// @Accept json
// @Produce json
// @Param input body product.CartDeleteRequest true "Данные для удаления корзины"
// @Success 200 {object} map[string]interface{} "Успешное удаление"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/ [delete]
func (r *Api) DeleteCartHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var request product.CartDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.logger.Error("Failed to binding the JSON", zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := r.product.DeleteCart(ctx, request)
	if err != nil {
		r.logger.Error("Failed to delete products", zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete cart")
		return
	}

	response.NewResponse(c, http.StatusOK, "cart deleted successfully", nil)
}

// UpdateCartHandler обновляет корзину.
// @Summary Обновить корзину
// @Description Обновляет корзину пользователя
// @Tags product
// @Accept json
// @Produce json
// @Param input body product.CartUpdateRequest true "Данные для обновления корзины"
// @Success 200 {object} map[string]interface{} "Успешное обновление"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/ [post]
func (r *Api) UpdateCartHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var request product.CartUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.logger.Error("Failed to binding the JSON", zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := r.product.UpdateCart(ctx, request)
	if err != nil {
		r.logger.Error("Failed to update cart", zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "failed to update cart")
		return
	}

	response.NewResponse(c, http.StatusOK, "cart updated successfully", nil)
}
