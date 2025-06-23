package serviceimplement

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type InventoryHistoryService struct {
	inventoryHistoryRepository repository.InventoryHistoryRepository
}

func NewInventoryHistoryService(inventoryHistoryRepository repository.InventoryHistoryRepository) service.InventoryHistoryService {
	return &InventoryHistoryService{
		inventoryHistoryRepository: inventoryHistoryRepository,
	}
}

func (s *InventoryHistoryService) GetAll(ctx *gin.Context, productID int) (*model.GetAllInventoryHistoriesResponse, string) {
	// Get all inventory histories
	inventoryHistories, err := s.inventoryHistoryRepository.GetAllByProductIDQuery(ctx, productID, nil)
	if err != nil {
		log.Error("InventoryHistoryService.GetAll Error when get inventory histories: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models
	inventoryHistoryResponses := make([]model.InventoryHistoryResponse, len(inventoryHistories))
	for i, inventoryHistory := range inventoryHistories {
		inventoryHistoryResponses[i] = model.InventoryHistoryResponse{
			ID:            inventoryHistory.ID,
			ProductID:     inventoryHistory.ProductID,
			Quantity:      inventoryHistory.Quantity,
			FinalQuantity: inventoryHistory.FinalQuantity,
			ImporterName:  inventoryHistory.ImporterName,
			ImportedAt:    inventoryHistory.ImportedAt,
			Note:          inventoryHistory.Note,
			ReferenceID:   inventoryHistory.ReferenceID,
		}
	}

	return &model.GetAllInventoryHistoriesResponse{
		InventoryHistories: inventoryHistoryResponses,
	}, ""
}

func (s *InventoryHistoryService) Create(ctx *gin.Context, request model.CreateInventoryHistoryRequest) (*model.InventoryHistoryResponse, string) {
	// Create inventory history entity
	inventoryHistory := &entity.InventoryHistory{
		ProductID:    request.ProductID,
		Quantity:     request.Quantity,
		ImporterName: request.ImporterName,
		ImportedAt:   time.Now(),
		Note:         request.Note,
	}

	// Save to database
	err := s.inventoryHistoryRepository.CreateCommand(ctx, inventoryHistory, nil)
	if err != nil {
		log.Error("InventoryHistoryService.Create Error when create inventory history: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response
	return &model.InventoryHistoryResponse{
		ID:            inventoryHistory.ID,
		ProductID:     inventoryHistory.ProductID,
		Quantity:      inventoryHistory.Quantity,
		FinalQuantity: inventoryHistory.FinalQuantity,
		ImporterName:  inventoryHistory.ImporterName,
		ImportedAt:    inventoryHistory.ImportedAt,
		Note:          inventoryHistory.Note,
		ReferenceID:   inventoryHistory.ReferenceID,
	}, ""
}
