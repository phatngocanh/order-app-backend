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
	productHandler *ProductHandler,
	inventoryHandler *InventoryHandler,
	inventoryHistoryHandler *InventoryHistoryHandler,
	customerHandler *CustomerHandler,
	orderHandler *OrderHandler,
	orderImageHandler *OrderImageHandler,
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
		products := v1.Group("/products")
		{
			products.POST("", authMiddleware.VerifyAccessToken, productHandler.Create)
			products.PUT("", authMiddleware.VerifyAccessToken, productHandler.Update)
			products.GET("", authMiddleware.VerifyAccessToken, productHandler.GetAll)
			products.GET("/:productId", authMiddleware.VerifyAccessToken, productHandler.GetOne)
			products.GET("/:productId/inventories", authMiddleware.VerifyAccessToken, inventoryHandler.GetByProductID)
			products.PUT("/:productId/inventories/quantity", authMiddleware.VerifyAccessToken, inventoryHandler.UpdateQuantity)
			products.GET("/:productId/inventories/histories", authMiddleware.VerifyAccessToken, inventoryHistoryHandler.GetAll)
		}
		customers := v1.Group("/customers")
		{
			customers.POST("", authMiddleware.VerifyAccessToken, customerHandler.Create)
			customers.PUT("/:customerId", authMiddleware.VerifyAccessToken, customerHandler.Update)
			customers.GET("", authMiddleware.VerifyAccessToken, customerHandler.GetAll)
			customers.GET("/:customerId", authMiddleware.VerifyAccessToken, customerHandler.GetOne)
		}
		orders := v1.Group("/orders")
		{
			orders.POST("", authMiddleware.VerifyAccessToken, orderHandler.Create)
			orders.PUT("/:orderId", authMiddleware.VerifyAccessToken, orderHandler.Update)
			orders.GET("", authMiddleware.VerifyAccessToken, orderHandler.GetAll)
			orders.GET("/:orderId", authMiddleware.VerifyAccessToken, orderHandler.GetOne)
			orders.DELETE("/:orderId", authMiddleware.VerifyAccessToken, orderHandler.Delete)

			// Order images endpoints
			orders.POST("/:orderId/images", authMiddleware.VerifyAccessToken, orderImageHandler.UploadImage)
			orders.GET("/:orderId/images", authMiddleware.VerifyAccessToken, orderImageHandler.GetImagesByOrderID)
			orders.DELETE("/:orderId/images/:imageId", authMiddleware.VerifyAccessToken, orderImageHandler.DeleteImage)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
