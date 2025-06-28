package serviceimplement

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pna/order-app-backend/internal/controller/http/middleware"
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
	userRepository             repository.UserRepository
	productRepository          repository.ProductRepository
	unitOfWork                 repository.UnitOfWork
}

func NewInventoryService(
	inventoryRepository repository.InventoryRepository,
	inventoryHistoryRepository repository.InventoryHistoryRepository,
	userRepository repository.UserRepository,
	productRepository repository.ProductRepository,
	unitOfWork repository.UnitOfWork,
) service.InventoryService {
	return &InventoryService{
		inventoryRepository:        inventoryRepository,
		inventoryHistoryRepository: inventoryHistoryRepository,
		userRepository:             userRepository,
		productRepository:          productRepository,
		unitOfWork:                 unitOfWork,
	}
}

func (s *InventoryService) GetAll(ctx context.Context) (*model.GetAllInventoryResponse, string) {
	// Get all inventory
	inventories, err := s.inventoryRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("InventoryService.GetAll Error when get inventories: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models with product info
	inventoryResponses := make([]model.InventoryWithProductResponse, len(inventories))
	for i, inventory := range inventories {
		// Get product info for this inventory
		product, err := s.productRepository.GetOneByIDQuery(ctx, inventory.ProductID, nil)
		if err != nil {
			log.Error("InventoryService.GetAll Error when get product for inventory " + string(rune(inventory.ID)) + ": " + err.Error())
			// Continue without product info for this inventory
			inventoryResponses[i] = model.InventoryWithProductResponse{
				ID:        inventory.ID,
				ProductID: inventory.ProductID,
				Quantity:  inventory.Quantity,
				Version:   inventory.Version,
				Product: model.ProductInfo{
					ID:            inventory.ProductID,
					Name:          "N/A",
					Spec:          0,
					OriginalPrice: 0,
				},
			}
			continue
		}

		inventoryResponses[i] = model.InventoryWithProductResponse{
			ID:        inventory.ID,
			ProductID: inventory.ProductID,
			Quantity:  inventory.Quantity,
			Version:   inventory.Version,
			Product: model.ProductInfo{
				ID:            product.ID,
				Name:          product.Name,
				Spec:          product.Spec,
				OriginalPrice: product.OriginalPrice,
			},
		}
	}

	return &model.GetAllInventoryResponse{
		Inventories: inventoryResponses,
	}, ""
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
	// Get user ID from context
	userID := middleware.GetUserIdHelper(ctx)
	if userID == 0 {
		log.Error("InventoryService.UpdateQuantity Error: user ID not found in context")
		return nil, error_utils.ErrorCode.UNAUTHORIZED
	}

	// Get user details to get username
	user, err := s.userRepository.FindByIDQuery(ctx, int(userID), nil)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when get user: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if user == nil {
		log.Error("InventoryService.UpdateQuantity Error: user not found")
		return nil, error_utils.ErrorCode.UNAUTHORIZED
	}

	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
			log.Error("InventoryService.UpdateQuantity Error when rollback transaction: " + rollbackErr.Error())
		}
	}()

	toBeLockedInventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, productID, tx)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when get inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if toBeLockedInventory == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Get inventory with FOR UPDATE lock
	existingInventory, err := s.inventoryRepository.GetOneByIDForUpdateQuery(ctx, toBeLockedInventory.ID, tx)
	if err != nil {
		log.Error("InventoryService.UpdateQuantity Error when get inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingInventory == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Check if version matches
	if existingInventory.Version != request.Version {
		return nil, error_utils.ErrorCode.INVENTORY_VERSION_MISMATCH
	}

	// Generate new version UUID
	newVersion := uuid.New().String()

	// Update inventory quantity with version check
	err = s.inventoryRepository.UpdateQuantityWithVersionCommand(ctx, productID, request.Quantity, request.Version, newVersion, tx)
	if err != nil {
		// Check for specific error types
		var constraintViolationError *error_utils.ConstraintViolationError
		if errors.As(err, &constraintViolationError) {
			return nil, error_utils.ErrorCode.INVENTORY_QUANTITY_NEGATIVE
		}

		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Create inventory history record
	inventoryHistory := &entity.InventoryHistory{
		ProductID:     productID,
		Quantity:      request.Quantity,
		FinalQuantity: existingInventory.Quantity + request.Quantity,
		ImporterName:  user.Username,
		ImportedAt:    time.Now(),
		Note:          request.Note,
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
