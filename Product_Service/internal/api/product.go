package api

import (
	"Product_Service/internal/api/response"
	"Product_Service/internal/domain"
	"Product_Service/internal/erorrs"
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// GetAllProducts возвращает список всех продуктов.
// @Summary Получить все продукты
// @Description РАБОТАЕТ ТОЛЬКО ОТ СЕРВИСА ЮЗЕРА, САМ ПО СЕБЕ НЕ РАБОТАЕТ
// @Tags products
// @Produce json
// @Success 200 {object} map[string]interface{} "Список продуктов"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/ [get]
func (r *ProductApi) GetAllProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "0"))
	if err != nil {
		r.logger.Error("Invalid page_size parameter",
			zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid page_size")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		r.logger.Error("Invalid page parameter",
			zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid page")
		return
	}

	products, err := r.service.GetAllProducts(ctx, pageSize, page)
	if err != nil {
		r.logger.Error("Error getting products",
			zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, products)
}

// CreateCartItem создает элемент корзины.
// @Summary Создать элемент корзины
// @Description Создает новый элемент корзины для пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Param input body domain.CartDTO true "Данные для создания корзины"
// @Success 200 {object} map[string]interface{} "Успешное создание"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/cart [post]
func (r *ProductApi) CreateCartItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var input domain.CartDTO

	err := c.BindJSON(&input)
	if err != nil {
		r.logger.Error("Error binding JSON",
			zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "bad request")
		return
	}

	err = r.service.CreateCart(ctx, &input)
	if err != nil {
		r.logger.Error("Error create cart",
			zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "success", nil)

}

// UpdateCartItem обновляет элемент корзины.
// @Summary Обновить элемент корзины
// @Description Обновляет существующий элемент корзины
// @Tags cart
// @Accept json
// @Produce json
// @Param input body domain.CartDTO true "Данные для обновления корзины"
// @Success 200 {object} map[string]interface{} "Успешное обновление"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/cart [put]
func (r *ProductApi) UpdateCartItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var input domain.CartDTO

	err := c.BindJSON(&input)
	if err != nil {
		r.logger.Error("Error binding JSON",
			zap.Error(err))
		response.NewErrorResponse(c, http.StatusBadRequest, "bad request")
		return
	}

	err = r.service.UpdateCartItem(ctx, &input)
	if err != nil {
		r.logger.Error("Error getting products",
			zap.Error(err))
		if errors.Is(err, erorrs.ErrCartNotFound) {
			r.logger.Error("Cart not found",
				zap.Error(err))
			response.NewErrorResponse(c, http.StatusBadRequest, "invalid user")
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "success", nil)

}

// DeleteCartItem удаляет элемент корзины.
// @Summary Удалить элемент корзины
// @Description Удаляет элемент корзины для пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Param input body domain.CartDeleteDTO true "Данные для удаления корзины"
// @Success 200 {object} map[string]interface{} "Успешное удаление"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/cart [delete]
func (r *ProductApi) DeleteCartItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var input domain.CartDeleteDTO

	err := c.BindJSON(&input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "internal error")
		return
	}

	err = r.service.DeleteCartItem(ctx, input.UserId, input.ProductId)
	if err != nil {
		r.logger.Error("Error deleting products",
			zap.Error(err))
		if errors.Is(err, erorrs.ErrCartNotFound) {
			r.logger.Error("Cart not found",
				zap.Error(err))
			response.NewErrorResponse(c, http.StatusBadRequest, "invalid user")
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	response.NewResponse(c, http.StatusOK, "success", nil)
}
