package serviceimplement

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type ProductService struct {
	productRepository   repository.ProductRepository
	inventoryRepository repository.InventoryRepository
	unitOfWork          repository.UnitOfWork
}

func NewProductService(productRepository repository.ProductRepository, inventoryRepository repository.InventoryRepository, unitOfWork repository.UnitOfWork) service.ProductService {
	return &ProductService{
		productRepository:   productRepository,
		inventoryRepository: inventoryRepository,
		unitOfWork:          unitOfWork,
	}
}

func (s *ProductService) Create(ctx *gin.Context, request model.CreateProductRequest) (*model.ProductResponse, string) {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("ProductService.Create Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("ProductService.Create Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Create product entity
	product := &entity.Product{
		Name:          request.Name,
		Spec:          request.Spec,
		OriginalPrice: request.OriginalPrice,
	}

	// Save product to database
	err = s.productRepository.CreateCommand(ctx, product, tx)
	if err != nil {
		log.Error("ProductService.Create Error when create product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Create inventory for the product
	inventory := &entity.Inventory{
		ProductID: product.ID,
		Quantity:  0, // Start with 0 quantity
		Version:   uuid.New().String(),
	}

	err = s.inventoryRepository.CreateCommand(ctx, inventory, tx)
	if err != nil {
		log.Error("ProductService.Create Error when create inventory: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("ProductService.Create Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response with inventory info
	return &model.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Spec:          product.Spec,
		OriginalPrice: product.OriginalPrice,
		Inventory: &model.InventoryInfo{
			Quantity: inventory.Quantity,
			Version:  inventory.Version,
		},
	}, ""
}

func (s *ProductService) Update(ctx *gin.Context, request model.UpdateProductRequest) (*model.ProductResponse, string) {
	// Check if product exists
	existingProduct, err := s.productRepository.GetOneByIDQuery(ctx, request.ID, nil)
	if err != nil {
		log.Error("ProductService.Update Error when get product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingProduct == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Update product entity
	product := &entity.Product{
		ID:            request.ID,
		Name:          request.Name,
		Spec:          request.Spec,
		OriginalPrice: request.OriginalPrice,
	}

	// Save to database
	err = s.productRepository.UpdateCommand(ctx, product, nil)
	if err != nil {
		log.Error("ProductService.Update Error when update product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Get inventory info for response
	inventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, product.ID, nil)
	if err != nil {
		log.Error("ProductService.Update Error when get inventory: " + err.Error())
		// Don't fail the update, just return without inventory info
		return &model.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			Spec:          product.Spec,
			OriginalPrice: product.OriginalPrice,
		}, ""
	}

	// Return response with inventory info
	return &model.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Spec:          product.Spec,
		OriginalPrice: product.OriginalPrice,
		Inventory: &model.InventoryInfo{
			Quantity: inventory.Quantity,
			Version:  inventory.Version,
		},
	}, ""
}

func (s *ProductService) GetAll(ctx *gin.Context) (*model.GetAllProductsResponse, string) {
	// Get all products
	products, err := s.productRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("ProductService.GetAll Error when get products: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models with inventory info
	productResponses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		// Get inventory for this product
		inventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, product.ID, nil)
		if err != nil {
			log.Error("ProductService.GetAll Error when get inventory for product " + string(rune(product.ID)) + ": " + err.Error())
			// Continue without inventory info for this product
			productResponses[i] = model.ProductResponse{
				ID:            product.ID,
				Name:          product.Name,
				Spec:          product.Spec,
				OriginalPrice: product.OriginalPrice,
			}
			continue
		}

		productResponses[i] = model.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			Spec:          product.Spec,
			OriginalPrice: product.OriginalPrice,
			Inventory: &model.InventoryInfo{
				Quantity: inventory.Quantity,
				Version:  inventory.Version,
			},
		}
	}

	return &model.GetAllProductsResponse{
		Products: productResponses,
	}, ""
}

func (s *ProductService) GetOne(ctx *gin.Context, id int) (*model.GetOneProductResponse, string) {
	// Get product by ID
	product, err := s.productRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("ProductService.GetOne Error when get product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if product == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Get inventory for this product
	inventory, err := s.inventoryRepository.GetOneByProductIDQuery(ctx, product.ID, nil)
	if err != nil {
		log.Error("ProductService.GetOne Error when get inventory: " + err.Error())
		// Return product without inventory info
		return &model.GetOneProductResponse{
			Product: model.ProductResponse{
				ID:            product.ID,
				Name:          product.Name,
				Spec:          product.Spec,
				OriginalPrice: product.OriginalPrice,
			},
		}, ""
	}

	// Return response with inventory info
	return &model.GetOneProductResponse{
		Product: model.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			Spec:          product.Spec,
			OriginalPrice: product.OriginalPrice,
			Inventory: &model.InventoryInfo{
				Quantity: inventory.Quantity,
				Version:  inventory.Version,
			},
		},
	}, ""
}
