package serviceimplement

import (
	"context"

	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type StatisticsService struct {
	productRepo   repository.ProductRepository
	customerRepo  repository.CustomerRepository
	inventoryRepo repository.InventoryRepository
	orderRepo     repository.OrderRepository
}

func NewStatisticsService(
	productRepo repository.ProductRepository,
	customerRepo repository.CustomerRepository,
	inventoryRepo repository.InventoryRepository,
	orderRepo repository.OrderRepository,
) service.StatisticsService {
	return &StatisticsService{
		productRepo:   productRepo,
		customerRepo:  customerRepo,
		inventoryRepo: inventoryRepo,
		orderRepo:     orderRepo,
	}
}

func (s *StatisticsService) GetDashboardStats(ctx context.Context) (model.DashboardStatsResponse, string) {
	// Get total products
	products, err := s.productRepo.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("StatisticsService.GetDashboardStats Error fetching products: " + err.Error())
		return model.DashboardStatsResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	// Get total customers
	customers, err := s.customerRepo.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("StatisticsService.GetDashboardStats Error fetching customers: " + err.Error())
		return model.DashboardStatsResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	// Get total orders
	orders, err := s.orderRepo.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("StatisticsService.GetDashboardStats Error fetching orders: " + err.Error())
		return model.DashboardStatsResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	// Calculate inventory stats
	totalInventoryItems := 0
	lowStockProducts := 0

	for _, product := range products {
		inventory, err := s.inventoryRepo.GetOneByProductIDQuery(ctx, product.ID, nil)
		if err != nil {
			log.Error("StatisticsService.GetDashboardStats Error fetching inventory for product " + string(rune(product.ID)) + ": " + err.Error())
			continue
		}
		if inventory != nil {
			totalInventoryItems += inventory.Quantity
			if inventory.Quantity < 10 { // Consider low stock if less than 10
				lowStockProducts++
			}
		}
	}

	// Calculate pending orders
	pendingOrders := 0
	for _, order := range orders {
		if order.DeliveryStatus != "COMPLETED" {
			pendingOrders++
		}
	}

	return model.DashboardStatsResponse{
		TotalProducts:       len(products),
		TotalCustomers:      len(customers),
		TotalInventoryItems: totalInventoryItems,
		LowStockProducts:    lowStockProducts,
		TotalOrders:         len(orders),
		PendingOrders:       pendingOrders,
	}, ""
}
