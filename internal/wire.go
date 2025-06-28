//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	beanimplement "github.com/pna/order-app-backend/internal/bean/implement"
	"github.com/pna/order-app-backend/internal/controller"
	"github.com/pna/order-app-backend/internal/controller/http"
	"github.com/pna/order-app-backend/internal/controller/http/middleware"
	v1 "github.com/pna/order-app-backend/internal/controller/http/v1"
	"github.com/pna/order-app-backend/internal/database"
	repositoryimplement "github.com/pna/order-app-backend/internal/repository/implement"
	serviceimplement "github.com/pna/order-app-backend/internal/service/implement"
)

var container = wire.NewSet(
	controller.NewApiContainer,
)

// may have grpc server in the future
var serverSet = wire.NewSet(
	http.NewServer,
)

// handler === controller | with service and repository layers to form 3 layers architecture
var handlerSet = wire.NewSet(
	v1.NewHealthHandler,
	v1.NewHelloWorldHandler,
	v1.NewUserHandler,
	v1.NewProductHandler,
	v1.NewInventoryHandler,
	v1.NewInventoryHistoryHandler,
	v1.NewCustomerHandler,
	v1.NewOrderHandler,
	v1.NewOrderImageHandler,
	v1.NewStatisticsHandler,
)

var serviceSet = wire.NewSet(
	serviceimplement.NewHelloWorldService,
	serviceimplement.NewUserService,
	serviceimplement.NewProductService,
	serviceimplement.NewInventoryService,
	serviceimplement.NewInventoryHistoryService,
	serviceimplement.NewCustomerService,
	serviceimplement.NewOrderService,
	serviceimplement.NewOrderImageService,
	serviceimplement.NewStatisticsService,
)

var repositorySet = wire.NewSet(
	repositoryimplement.NewHelloWorldRepository,
	repositoryimplement.NewUserRepository,
	repositoryimplement.NewProductRepository,
	repositoryimplement.NewInventoryRepository,
	repositoryimplement.NewInventoryHistoryRepository,
	repositoryimplement.NewUnitOfWork,
	repositoryimplement.NewCustomerRepository,
	repositoryimplement.NewOrderRepository,
	repositoryimplement.NewOrderItemRepository,
	repositoryimplement.NewOrderImageRepository,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
)

var beanSet = wire.NewSet(
	beanimplement.NewBcryptPasswordEncoder,
	beanimplement.NewS3Service,
)

func InitializeContainer(
	db database.Db,
) *controller.ApiContainer {
	wire.Build(serverSet, handlerSet, serviceSet, repositorySet, middlewareSet, beanSet, container)
	return &controller.ApiContainer{}
}
