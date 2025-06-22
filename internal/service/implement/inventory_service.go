package serviceimplement

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type InventoryService struct {
	inventoryRepository repository.InventoryRepository
}

func NewInventoryService(inventoryRepository repository.InventoryRepository) service.InventoryService {
	return &InventoryService{
		inventoryRepository: inventoryRepository,
	}
}

func (s *InventoryService) GetByProductID(ctx *gin.Context, productID int) (*model.InventoryResponse, string) {
	// Get inventory by product ID
	inventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, productID, nil)
	if err != nil {
		log.Error("InventoryService.GetByProductID Error when get inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if inventory == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Return response
	return &model.InventoryResponse{
		ID:        inventory.ID,
		ProductID: inventory.ProductID,
		Quantity:  inventory.Quantity,
		Version:   inventory.Version,
	}, ""
}

func (s *InventoryService) UpdateQuantity(ctx *gin.Context, request model.UpdateInventoryQuantityRequest) (*model.InventoryResponse, string) {
	// Check if inventory exists
	existingInventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, request.ProductID, nil)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when get inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingInventory == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Generate new version UUID
	newVersion := uuid.New().String()

	// Update inventory quantity
	err = s.inventoryRepository.UpdateQuantityCommand(ctx, request.ProductID, request.Quantity, newVersion, nil)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when update inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Get updated inventory
	updatedInventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, request.ProductID, nil)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when get updated inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response
	return &model.InventoryResponse{
		ID:        updatedInventory.ID,
		ProductID: updatedInventory.ProductID,
		Quantity:  updatedInventory.Quantity,
		Version:   updatedInventory.Version,
	}, ""
}
