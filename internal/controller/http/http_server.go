package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/pna/order-app-backend/internal/controller/http/middleware"

	"github.com/gin-gonic/gin"

	v1 "github.com/pna/order-app-backend/internal/controller/http/v1"
)

type Server struct {
	healthHandler           *v1.HealthHandler
	helloWorldHandler       *v1.HelloWorldHandler
	authMiddleware          *middleware.AuthMiddleware
	userHandler             *v1.UserHandler
	productHandler          *v1.ProductHandler
	inventoryHandler        *v1.InventoryHandler
	inventoryHistoryHandler *v1.InventoryHistoryHandler
	customerHandler         *v1.CustomerHandler
	orderHandler            *v1.OrderHandler
	orderImageHandler       *v1.OrderImageHandler
	statisticsHandler       *v1.StatisticsHandler
}

func NewServer(
	healthHandler *v1.HealthHandler,
	helloWorldHandler *v1.HelloWorldHandler,
	authMiddleware *middleware.AuthMiddleware,
	userHandler *v1.UserHandler,
	productHandler *v1.ProductHandler,
	inventoryHandler *v1.InventoryHandler,
	inventoryHistoryHandler *v1.InventoryHistoryHandler,
	customerHandler *v1.CustomerHandler,
	orderHandler *v1.OrderHandler,
	orderImageHandler *v1.OrderImageHandler,
	statisticsHandler *v1.StatisticsHandler,
) *Server {
	return &Server{
		healthHandler:           healthHandler,
		helloWorldHandler:       helloWorldHandler,
		authMiddleware:          authMiddleware,
		userHandler:             userHandler,
		productHandler:          productHandler,
		inventoryHandler:        inventoryHandler,
		inventoryHistoryHandler: inventoryHistoryHandler,
		customerHandler:         customerHandler,
		orderHandler:            orderHandler,
		orderImageHandler:       orderImageHandler,
		statisticsHandler:       statisticsHandler,
	}
}

func (s *Server) Run() {
	router := gin.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	httpServerInstance := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	fmt.Println("Server running at " + httpServerInstance.Addr)

	v1.MapRoutes(
		router,
		s.healthHandler,
		s.helloWorldHandler,
		s.userHandler,
		s.productHandler,
		s.inventoryHandler,
		s.inventoryHistoryHandler,
		s.customerHandler,
		s.orderHandler,
		s.orderImageHandler,
		s.statisticsHandler,
		s.authMiddleware,
	)
	err := httpServerInstance.ListenAndServe()
	if err != nil {
		fmt.Println("There is error: " + err.Error())
		return
	}
}
