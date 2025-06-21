package serviceimplement

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type ProductService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) service.ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) Create(ctx *gin.Context, request model.CreateProductRequest) (*model.ProductResponse, string) {
	// Create product entity
	product := &entity.Product{
		Name:          request.Name,
		Spec:          request.Spec,
		OriginalPrice: request.OriginalPrice,
	}

	// Save to database
	err := s.productRepository.CreateCommand(ctx, product, nil)
	if err != nil {
		log.Error("ProductService.Create Error when create product: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response
	return &model.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Spec:          product.Spec,
		OriginalPrice: product.OriginalPrice,
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

	// Return response
	return &model.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Spec:          product.Spec,
		OriginalPrice: product.OriginalPrice,
	}, ""
}

func (s *ProductService) GetAll(ctx *gin.Context) (*model.GetAllProductsResponse, string) {
	// Get all products
	products, err := s.productRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("ProductService.GetAll Error when get products: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models
	productResponses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = model.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			Spec:          product.Spec,
			OriginalPrice: product.OriginalPrice,
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

	// Return response
	return &model.GetOneProductResponse{
		Product: model.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			Spec:          product.Spec,
			OriginalPrice: product.OriginalPrice,
		},
	}, ""
}
