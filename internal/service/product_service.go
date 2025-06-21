package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/model"
)

type ProductService interface {
	Create(ctx *gin.Context, request model.CreateProductRequest) (*model.ProductResponse, string)
	Update(ctx *gin.Context, request model.UpdateProductRequest) (*model.ProductResponse, string)
	GetAll(ctx *gin.Context) (*model.GetAllProductsResponse, string)
	GetOne(ctx *gin.Context, id int) (*model.GetOneProductResponse, string)
}
