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
	userHandler *UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Apply CORS middleware to all routes
	router.Use(middleware.CorsMiddleware())

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
		users := v1.Group("/users")
		{
			users.POST("/login", userHandler.Login)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
