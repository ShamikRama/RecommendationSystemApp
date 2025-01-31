package api

import (
	"User_Service/internal/logger"
	"User_Service/internal/service"
	"User_Service/pkg/product"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Api struct {
	service *service.Service
	product product.Client
	logger  logger.Logger
}

func NewApi(service *service.Service, product product.Client, logger logger.Logger) *Api {
	return &Api{
		service: service,
		product: product,
		logger:  logger,
	}
}

func (r *Api) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", r.SignUp)
		auth.POST("/sign-in", r.SignIn)
	}

	api := router.Group("/users", r.UserIdentity)
	{
		api.PUT("/:id", r.UpdateUser)
	}

	products := router.Group("/products")
	{
		products.GET("", r.GetProductsHandler)    // all products
		products.POST("/", r.CreateCartHandler)   // post for user
		products.DELETE("/", r.DeleteCartHandler) // delete product for user
		products.PUT("/", r.UpdateCartHandler)    // put  product for user
	}

	return router
}
