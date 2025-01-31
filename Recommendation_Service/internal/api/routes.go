package api

import (
	"Recommendation_Service/internal/logger"
	"Recommendation_Service/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type RecApi struct {
	service service.ServiceRec
	logger  logger.Logger
}

func NewRecApi(service service.ServiceRec, logger logger.Logger) *RecApi {
	return &RecApi{
		service: service,
		logger:  logger,
	}
}

func (r *RecApi) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/rec/{id}", r.GetRecommend)

	return router
}
