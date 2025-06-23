package serviceimplement

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type InventoryService struct {
	inventoryRepository        repository.InventoryRepository
	inventoryHistoryRepository repository.InventoryHistoryRepository
	unitOfWork                 repository.UnitOfWork
}

func NewInventoryService(inventoryRepository repository.InventoryRepository, inventoryHistoryRepository repository.InventoryHistoryRepository, unitOfWork repository.UnitOfWork) service.InventoryService {
	return &InventoryService{
		inventoryRepository:        inventoryRepository,
		inventoryHistoryRepository: inventoryHistoryRepository,
		unitOfWork:                 unitOfWork,
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

func (s *InventoryService) UpdateQuantity(ctx *gin.Context, productID int, request model.UpdateInventoryQuantityRequest) (*model.InventoryResponse, string) {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("InventoryService.UpdateQuantity Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Check if inventory exists
	existingInventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, productID, tx)
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
	err = s.inventoryRepository.UpdateQuantityCommand(ctx, productID, request.Quantity, newVersion, tx)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when update inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Create inventory history record
	inventoryHistory := &entity.InventoryHistory{
		ProductID:    productID,
		Quantity:     request.Quantity,
		ImporterName: request.ImporterName,
		ImportedAt:   time.Now(),
		Note:         request.Note,
	}

	err = s.inventoryHistoryRepository.CreateCommand(ctx, inventoryHistory, tx)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when create inventory history: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Get updated inventory
	updatedInventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, productID, nil)
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
