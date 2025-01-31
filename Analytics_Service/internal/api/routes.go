package api

import (
	"Analytics_Service/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Api struct {
	service service.ServiceClient
}

func NewApi(service service.ServiceClient) *Api {
	return &Api{
		service: service,
	}
}

func InitRoutes(r *Api) *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	stat := router.Group("/stats")
	{
		stat.GET("/products", r.GetStats)
		stat.GET("/actions", r.GetUserActions)
		stat.GET("/events", r.GetUserEvents)
	}

	return router

}
