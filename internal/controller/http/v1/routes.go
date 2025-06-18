package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/controller/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapRoutes(router *gin.Engine,
	healHandler *HealthHandler,
	helloWorldHandler *HelloWorldHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	v1 := router.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			health.GET("", healHandler.Check)
		}
		hello := v1.Group("/hello-world")
		{
			hello.GET("", helloWorldHandler.HelloWorld)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
