package http

import (
	"fmt"
	"github.com/pna/order-app-backend/internal/controller/http/middleware"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	v1 "github.com/pna/order-app-backend/internal/controller/http/v1"
)

type Server struct {
	healthHandler     *v1.HealthHandler
	helloWorldHandler *v1.HelloWorldHandler
	authMiddleware    *middleware.AuthMiddleware
}

func NewServer(
	healthHandler *v1.HealthHandler,
	helloWorldHandler *v1.HelloWorldHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Server {
	return &Server{
		healthHandler:     healthHandler,
		helloWorldHandler: helloWorldHandler,
		authMiddleware:    authMiddleware,
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
		s.authMiddleware,
	)
	err := httpServerInstance.ListenAndServe()
	if err != nil {
		fmt.Println("There is error: " + err.Error())
		return
	}
}
