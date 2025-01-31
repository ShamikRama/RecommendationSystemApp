package api

import (
	"Product_Service/internal/logger"
	"Product_Service/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type ProductApi struct {
	service service.ServiceProduct
	logger  logger.Logger
}

func NewProductApi(service service.ServiceProduct, logger logger.Logger) *ProductApi {
	return &ProductApi{
		service: service,
		logger:  logger,
	}
}

func (r *ProductApi) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	product := router.Group("/products") // было products
	product.Use(AuthMiddleware())
	{
		product.GET("", r.GetAllProducts)
		product.POST("/cart", r.CreateCartItem)
		product.PUT("/cart", r.UpdateCartItem)
		product.DELETE("/cart", r.DeleteCartItem)
	}

	return router
}
