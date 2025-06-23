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

type CustomerService struct {
	customerRepository repository.CustomerRepository
	unitOfWork         repository.UnitOfWork
}

func NewCustomerService(customerRepository repository.CustomerRepository, unitOfWork repository.UnitOfWork) service.CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
		unitOfWork:         unitOfWork,
	}
}

func (s *CustomerService) Create(ctx *gin.Context, request model.CreateCustomerRequest) (*model.CustomerResponse, string) {
	// Begin transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("CustomerService.Create Error when begin transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("CustomerService.Create Error when rollback transaction: " + rollbackErr.Error())
			}
		}
	}()

	// Create customer entity
	customer := &entity.Customer{
		Name:    request.Name,
		Phone:   request.Phone,
		Address: request.Address,
	}

	// Save customer to database
	err = s.customerRepository.CreateCommand(ctx, customer, tx)
	if err != nil {
		log.Error("CustomerService.Create Error when create customer: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("CustomerService.Create Error when commit transaction: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response
	return &model.CustomerResponse{
		ID:      customer.ID,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Address: customer.Address,
	}, ""
}

func (s *CustomerService) Update(ctx *gin.Context, customerID int, request model.UpdateCustomerRequest) (*model.CustomerResponse, string) {
	// Check if customer exists
	existingCustomer, err := s.customerRepository.GetOneByIDQuery(ctx, customerID, nil)
	if err != nil {
		log.Error("CustomerService.Update Error when get customer: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if existingCustomer == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Update customer entity - only update non-empty fields
	customer := &entity.Customer{
		ID:      customerID,
		Name:    existingCustomer.Name,    // Keep existing value
		Phone:   existingCustomer.Phone,   // Keep existing value
		Address: existingCustomer.Address, // Keep existing value
	}

	// Only update fields that are not empty
	if request.Name != "" {
		customer.Name = request.Name
	}
	if request.Phone != "" {
		customer.Phone = request.Phone
	}
	if request.Address != "" {
		customer.Address = request.Address
	}

	// Save to database
	err = s.customerRepository.UpdateCommand(ctx, customer, nil)
	if err != nil {
		log.Error("CustomerService.Update Error when update customer: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Return response
	return &model.CustomerResponse{
		ID:      customer.ID,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Address: customer.Address,
	}, ""
}

func (s *CustomerService) GetAll(ctx *gin.Context) (*model.GetAllCustomersResponse, string) {
	// Get all customers
	customers, err := s.customerRepository.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("CustomerService.GetAll Error when get customers: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	// Convert to response models
	customerResponses := make([]model.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = model.CustomerResponse{
			ID:      customer.ID,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		}
	}

	return &model.GetAllCustomersResponse{
		Customers: customerResponses,
	}, ""
}

func (s *CustomerService) GetOne(ctx *gin.Context, id int) (*model.GetOneCustomerResponse, string) {
	// Get customer by ID
	customer, err := s.customerRepository.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("CustomerService.GetOne Error when get customer: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	if customer == nil {
		return nil, error_utils.ErrorCode.NOT_FOUND
	}

	// Return response
	return &model.GetOneCustomerResponse{
		Customer: model.CustomerResponse{
			ID:      customer.ID,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		},
	}, ""
}
