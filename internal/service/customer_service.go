package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/domain/model"
)

type CustomerService interface {
	Create(ctx *gin.Context, request model.CreateCustomerRequest) (*model.CustomerResponse, string)
	Update(ctx *gin.Context, customerID int, request model.UpdateCustomerRequest) (*model.CustomerResponse, string)
	GetAll(ctx *gin.Context) (*model.GetAllCustomersResponse, string)
	GetOne(ctx *gin.Context, id int) (*model.GetOneCustomerResponse, string)
}
