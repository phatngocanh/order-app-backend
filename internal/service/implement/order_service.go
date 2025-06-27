package serviceimplement

import (
	"context"
	"errors"
	"strconv"
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

type OrderService struct {
	orderRepo            repository.OrderRepository
	orderItemRepo        repository.OrderItemRepository
	inventoryRepo        repository.InventoryRepository
	inventoryHistoryRepo repository.InventoryHistoryRepository
	userRepo             repository.UserRepository
	productRepo          repository.ProductRepository
	unitOfWork           repository.UnitOfWork
	customerRepo         repository.CustomerRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	inventoryRepo repository.InventoryRepository,
	inventoryHistoryRepo repository.InventoryHistoryRepository,
	userRepo repository.UserRepository,
	unitOfWork repository.UnitOfWork,
	productRepo repository.ProductRepository,
	customerRepo repository.CustomerRepository) service.OrderService {
	return &OrderService{
		orderRepo:            orderRepo,
		orderItemRepo:        orderItemRepo,
		inventoryRepo:        inventoryRepo,
		inventoryHistoryRepo: inventoryHistoryRepo,
		userRepo:             userRepo,
		unitOfWork:           unitOfWork,
		productRepo:          productRepo,
		customerRepo:         customerRepo,
	}
}

// Helper to calculate total amount and product count from order items
func calculateOrderAmountsAndProductCount(orderItems []entity.OrderItem) (totalAmount int, productCount int) {
	productIDSet := make(map[int]struct{})
	totalAmount = 0
	for _, item := range orderItems {
		finalAmount := item.FinalAmount
		if finalAmount == nil {
			itemTotal := item.Quantity * item.SellingPrice
			discountAmount := (itemTotal * item.Discount) / 100
			calculated := itemTotal - discountAmount
			finalAmount = &calculated
		}
		totalAmount += *finalAmount
		productIDSet[item.ProductID] = struct{}{}
	}
	productCount = len(productIDSet)
	return
}

func (s *OrderService) GetAll(ctx context.Context) (model.GetAllOrdersResponse, string) {
	orders, err := s.orderRepo.GetAllQuery(ctx, nil)
	if err != nil {
		log.Error("OrderService.GetAll Error: " + err.Error())
		return model.GetAllOrdersResponse{}, error_utils.ErrorCode.DB_DOWN
	}
	resp := model.GetAllOrdersResponse{Orders: make([]model.OrderResponse, 0, len(orders))}
	for _, o := range orders {
		// Fetch customer information
		customer, err := s.customerRepo.GetOneByIDQuery(ctx, o.CustomerID, nil)
		if err != nil {
			log.Error("OrderService.GetAll Error fetching customer: " + err.Error())
			continue
		}

		// Fetch order items to calculate total amount and product count
		orderItems, err := s.orderItemRepo.GetAllByOrderIDQuery(ctx, o.ID, nil)
		if err != nil {
			log.Error("OrderService.GetAll Error fetching order items: " + err.Error())
			continue
		}
		totalAmount, productCount := calculateOrderAmountsAndProductCount(orderItems)

		resp.Orders = append(resp.Orders, model.OrderResponse{
			ID:                   o.ID,
			OrderDate:            o.OrderDate,
			DeliveryStatus:       o.DeliveryStatus,
			DebtStatus:           o.DebtStatus,
			StatusTransitionedAt: o.StatusTransitionedAt,
			ShippingFee:          o.ShippingFee,
			Customer: model.CustomerResponse{
				ID:      customer.ID,
				Name:    customer.Name,
				Phone:   customer.Phone,
				Address: customer.Address,
			},
			OrderItems:   nil, // Omit order items in GetAll
			TotalAmount:  &totalAmount,
			ProductCount: &productCount,
		})
	}
	return resp, ""
}

func (s *OrderService) GetOne(ctx context.Context, id int) (model.GetOneOrderResponse, string) {
	order, err := s.orderRepo.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("OrderService.GetOne Error: " + err.Error())
		return model.GetOneOrderResponse{}, error_utils.ErrorCode.DB_DOWN
	}
	if order == nil {
		return model.GetOneOrderResponse{}, error_utils.ErrorCode.NOT_FOUND
	}

	// Fetch customer information
	customer, err := s.customerRepo.GetOneByIDQuery(ctx, order.CustomerID, nil)
	if err != nil {
		log.Error("OrderService.GetOne Error fetching customer: " + err.Error())
		return model.GetOneOrderResponse{}, error_utils.ErrorCode.DB_DOWN
	}

	// Fetch order items
	orderItems, err := s.orderItemRepo.GetAllByOrderIDQuery(ctx, order.ID, nil)
	if err != nil {
		log.Error("OrderService.GetOne Error fetching order items: " + err.Error())
		return model.GetOneOrderResponse{}, error_utils.ErrorCode.DB_DOWN
	}
	orderItemResponses := make([]model.OrderItemResponse, 0, len(orderItems))
	for _, item := range orderItems {
		product, err := s.productRepo.GetOneByIDQuery(ctx, item.ProductID, nil)
		if err != nil {
			log.Error("OrderService.GetOne Error fetching product: " + err.Error())
			return model.GetOneOrderResponse{}, error_utils.ErrorCode.DB_DOWN
		}
		finalAmount := item.FinalAmount
		if finalAmount == nil {
			itemTotal := item.Quantity * item.SellingPrice
			discountAmount := (itemTotal * item.Discount) / 100
			calculated := itemTotal - discountAmount
			finalAmount = &calculated
		}
		orderItemResponses = append(orderItemResponses, model.OrderItemResponse{
			ID:            item.ID,
			OrderID:       item.OrderID,
			ProductID:     item.ProductID,
			ProductName:   product.Name,
			NumberOfBoxes: item.NumberOfBoxes,
			Spec:          item.Spec,
			Quantity:      item.Quantity,
			SellingPrice:  item.SellingPrice,
			Discount:      item.Discount,
			FinalAmount:   finalAmount,
			ExportFrom:    item.ExportFrom,
		})
	}
	totalAmount, productCount := calculateOrderAmountsAndProductCount(orderItems)

	resp := model.GetOneOrderResponse{Order: model.OrderResponse{
		ID:                   order.ID,
		OrderDate:            order.OrderDate,
		DeliveryStatus:       order.DeliveryStatus,
		DebtStatus:           order.DebtStatus,
		StatusTransitionedAt: order.StatusTransitionedAt,
		ShippingFee:          order.ShippingFee,
		Customer: model.CustomerResponse{
			ID:      customer.ID,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		},
		OrderItems:   orderItemResponses,
		TotalAmount:  &totalAmount,
		ProductCount: &productCount,
	}}
	return resp, ""
}

func (s *OrderService) Create(ctx *gin.Context, req model.CreateOrderRequest) string {
	// Get user ID from context
	userID := middleware.GetUserIdHelper(ctx)
	if userID == 0 {
		log.Error("OrderService.Create Error: user ID not found in context")
		return error_utils.ErrorCode.UNAUTHORIZED
	}

	// Get user details to get username
	user, err := s.userRepo.FindByIDQuery(ctx, int(userID), nil)
	if err != nil {
		log.Error("OrderService.Create Error when get user: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	if user == nil {
		log.Error("OrderService.Create Error: user not found")
		return error_utils.ErrorCode.UNAUTHORIZED
	}

	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("OrderService.Create Error when begin transaction: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	defer func() {
		if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
			log.Error("OrderService.Create Error when rollback transaction: " + rollbackErr.Error())
		}
	}()

	productIDs := make([]int, 0, len(req.OrderItems))
	for _, item := range req.OrderItems {
		productIDs = append(productIDs, item.ProductID)
	}

	inventoryIDs, err := s.inventoryRepo.GetInventoryIDsByProductIDsQuery(ctx, productIDs, tx)
	if err != nil {
		log.Error("OrderService.Create Error when get inventory IDs: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	lockedInventories, err := s.inventoryRepo.SelectManyForUpdate(ctx, inventoryIDs, tx)
	if err != nil {
		log.Error("OrderService.Create Error when lock inventories: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	inventoryMap := make(map[int]*entity.Inventory)
	for i := range lockedInventories {
		inv := &lockedInventories[i]
		inventoryMap[inv.ProductID] = inv
	}

	orderEntity := entity.Order{
		CustomerID:     req.CustomerID,
		OrderDate:      req.OrderDate,
		DeliveryStatus: req.DeliveryStatus,
		DebtStatus:     req.DebtStatus,
		ShippingFee:    req.ShippingFee,
	}
	now := time.Now()
	orderEntity.StatusTransitionedAt = &now

	err = s.orderRepo.CreateCommand(ctx, &orderEntity, tx)
	if err != nil {
		log.Error("OrderService.Create Error when create order: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	for _, item := range req.OrderItems {
		inv := inventoryMap[item.ProductID]
		quantityToExport := item.Quantity
		itemVersion := item.Version
		originalInventoryQuantity := inv.Quantity // Store original inventory quantity

		if inv.Version != itemVersion {
			log.Error("OrderService.Create Error: inventory version mismatch for productID ", item.ProductID)
			return error_utils.ErrorCode.INVENTORY_VERSION_MISMATCH
		}

		// Case 1: Order quantity <= inventory quantity
		if inv.Quantity >= quantityToExport {
			// Calculate final amount for inventory item
			itemTotal := quantityToExport * item.SellingPrice
			discountAmount := (itemTotal * item.Discount) / 100
			finalAmount := itemTotal - discountAmount

			itemEntity := entity.OrderItem{
				ProductID:     item.ProductID,
				NumberOfBoxes: item.NumberOfBoxes,
				Spec:          item.Spec,
				Quantity:      quantityToExport,
				SellingPrice:  item.SellingPrice,
				Discount:      item.Discount,
				FinalAmount:   &finalAmount,
				OrderID:       orderEntity.ID,
				ExportFrom:    entity.OrderExportFrom.INVENTORY,
			}

			newVersion := uuid.New().String()
			err = s.inventoryRepo.UpdateQuantityWithVersionCommand(ctx, inv.ProductID, -quantityToExport, inv.Version, newVersion, tx)
			if err != nil {
				var constraintViolationError *error_utils.ConstraintViolationError
				if errors.As(err, &constraintViolationError) {
					log.Error("OrderService.Create Error: inventory quantity negative for productID ", item.ProductID)
					return error_utils.ErrorCode.INVENTORY_QUANTITY_NEGATIVE
				}
				log.Error("OrderService.Create Error when update inventory: " + err.Error())
				return error_utils.ErrorCode.DB_DOWN
			}

			// Create inventory history record
			inventoryHistory := &entity.InventoryHistory{
				ProductID:     inv.ProductID,
				Quantity:      -quantityToExport,
				FinalQuantity: inv.Quantity - quantityToExport,
				ImporterName:  user.Username,
				ImportedAt:    time.Now(),
				Note:          "Hàng trừ cho hoá đơn ID: " + strconv.Itoa(orderEntity.ID),
				ReferenceID:   &orderEntity.ID,
			}
			err = s.inventoryHistoryRepo.CreateCommand(ctx, inventoryHistory, tx)
			if err != nil {
				log.Error("OrderService.Create Error when create inventory history: " + err.Error())
				return error_utils.ErrorCode.DB_DOWN
			}

			err = s.orderItemRepo.CreateCommand(ctx, &itemEntity, tx)
			if err != nil {
				log.Error("OrderService.Create Error when create order item: " + err.Error())
				return error_utils.ErrorCode.DB_DOWN
			}

			inv.Quantity -= quantityToExport
			inv.Version = newVersion
		} else {
			// Case 2: Order quantity > inventory quantity

			// Create order item from inventory (if inventory has stock)
			if inv.Quantity > 0 {
				// Calculate final amount for inventory item
				inventoryItemTotal := inv.Quantity * item.SellingPrice
				inventoryDiscountAmount := (inventoryItemTotal * item.Discount) / 100
				inventoryFinalAmount := inventoryItemTotal - inventoryDiscountAmount

				inventoryItemEntity := entity.OrderItem{
					ProductID:     item.ProductID,
					NumberOfBoxes: item.NumberOfBoxes,
					Spec:          item.Spec,
					Quantity:      inv.Quantity,
					SellingPrice:  item.SellingPrice,
					Discount:      item.Discount,
					FinalAmount:   &inventoryFinalAmount,
					OrderID:       orderEntity.ID,
					ExportFrom:    entity.OrderExportFrom.INVENTORY,
				}

				newVersion := uuid.New().String()
				err = s.inventoryRepo.UpdateQuantityWithVersionCommand(ctx, inv.ProductID, -inv.Quantity, inv.Version, newVersion, tx)
				if err != nil {
					var constraintViolationError *error_utils.ConstraintViolationError
					if errors.As(err, &constraintViolationError) {
						log.Error("OrderService.Create Error: inventory quantity negative for productID ", item.ProductID)
						return error_utils.ErrorCode.INVENTORY_QUANTITY_NEGATIVE
					}
					log.Error("OrderService.Create Error when update inventory: " + err.Error())
					return error_utils.ErrorCode.DB_DOWN
				}

				// Create inventory history record for inventory export
				inventoryHistory := &entity.InventoryHistory{
					ProductID:     inv.ProductID,
					Quantity:      -inv.Quantity,
					FinalQuantity: 0,
					ImporterName:  user.Username,
					ImportedAt:    time.Now(),
					Note:          "Hàng trừ cho hoá đơn ID: " + strconv.Itoa(orderEntity.ID),
					ReferenceID:   &orderEntity.ID,
				}
				err = s.inventoryHistoryRepo.CreateCommand(ctx, inventoryHistory, tx)
				if err != nil {
					log.Error("OrderService.Create Error when create inventory history: " + err.Error())
					return error_utils.ErrorCode.DB_DOWN
				}

				err = s.orderItemRepo.CreateCommand(ctx, &inventoryItemEntity, tx)
				if err != nil {
					log.Error("OrderService.Create Error when create inventory order item: " + err.Error())
					return error_utils.ErrorCode.DB_DOWN
				}

				inv.Quantity = 0
				inv.Version = newVersion
			}

			// Create order item from external for remaining quantity
			// Use original inventory quantity to calculate remaining quantity
			remainingQuantity := quantityToExport - originalInventoryQuantity
			if remainingQuantity > 0 {
				// Calculate final amount for external item
				externalItemTotal := remainingQuantity * item.SellingPrice
				externalDiscountAmount := (externalItemTotal * item.Discount) / 100
				externalFinalAmount := externalItemTotal - externalDiscountAmount

				externalItemEntity := entity.OrderItem{
					ProductID:     item.ProductID,
					NumberOfBoxes: item.NumberOfBoxes,
					Spec:          item.Spec,
					Quantity:      remainingQuantity,
					SellingPrice:  item.SellingPrice,
					Discount:      item.Discount,
					FinalAmount:   &externalFinalAmount,
					OrderID:       orderEntity.ID,
					ExportFrom:    entity.OrderExportFrom.EXTERNAL,
				}

				err = s.orderItemRepo.CreateCommand(ctx, &externalItemEntity, tx)
				if err != nil {
					log.Error("OrderService.Create Error when create external order item: " + err.Error())
					return error_utils.ErrorCode.DB_DOWN
				}
			}
		}
	}

	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("OrderService.Create Error when commit transaction: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	return ""
}

func (s *OrderService) Update(ctx context.Context, req model.UpdateOrderRequest) string {
	existing, err := s.orderRepo.GetOneByIDQuery(ctx, req.ID, nil)
	if err != nil {
		log.Error("OrderService.Update Error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if existing == nil {
		return error_utils.ErrorCode.NOT_FOUND
	}

	if req.CustomerID != 0 {
		existing.CustomerID = req.CustomerID
	}
	if !req.OrderDate.IsZero() {
		existing.OrderDate = req.OrderDate
	}
	if req.DeliveryStatus != "" {
		existing.DeliveryStatus = req.DeliveryStatus
		now := time.Now()
		existing.StatusTransitionedAt = &now
	}
	if req.DebtStatus != "" {
		existing.DebtStatus = req.DebtStatus
	}
	// Update shipping fee if provided
	existing.ShippingFee = req.ShippingFee

	err = s.orderRepo.UpdateCommand(ctx, existing, nil)
	if err != nil {
		log.Error("OrderService.Update Error when update order: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	return ""
}

func (s *OrderService) Delete(ctx *gin.Context, id int) string {
	// Get user ID from context
	userID := middleware.GetUserIdHelper(ctx)
	if userID == 0 {
		log.Error("OrderService.Delete Error: user ID not found in context")
		return error_utils.ErrorCode.UNAUTHORIZED
	}

	// Get user details to get username
	user, err := s.userRepo.FindByIDQuery(ctx, int(userID), nil)
	if err != nil {
		log.Error("OrderService.Delete Error when get user: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	if user == nil {
		log.Error("OrderService.Delete Error: user not found")
		return error_utils.ErrorCode.UNAUTHORIZED
	}

	// Check if order exists
	order, err := s.orderRepo.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("OrderService.Delete Error when get order: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if order == nil {
		return error_utils.ErrorCode.NOT_FOUND
	}

	// Get order items
	orderItems, err := s.orderItemRepo.GetAllByOrderIDQuery(ctx, id, nil)
	if err != nil {
		log.Error("OrderService.Delete Error when get order items: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// Filter items that were exported from inventory
	inventoryItems := make([]entity.OrderItem, 0)
	for _, item := range orderItems {
		if item.ExportFrom == entity.OrderExportFrom.INVENTORY {
			inventoryItems = append(inventoryItems, item)
		}
	}

	// If there are inventory items, we need to restore them
	if len(inventoryItems) > 0 {
		tx, err := s.unitOfWork.Begin(ctx)
		if err != nil {
			log.Error("OrderService.Delete Error when begin transaction: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}
		defer func() {
			if rollbackErr := s.unitOfWork.Rollback(tx); rollbackErr != nil {
				log.Error("OrderService.Delete Error when rollback transaction: " + rollbackErr.Error())
			}
		}()

		// Get product IDs for inventory items
		productIDs := make([]int, 0, len(inventoryItems))
		for _, item := range inventoryItems {
			productIDs = append(productIDs, item.ProductID)
		}

		// Get inventory IDs and lock inventories
		inventoryIDs, err := s.inventoryRepo.GetInventoryIDsByProductIDsQuery(ctx, productIDs, tx)
		if err != nil {
			log.Error("OrderService.Delete Error when get inventory IDs: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}

		lockedInventories, err := s.inventoryRepo.SelectManyForUpdate(ctx, inventoryIDs, tx)
		if err != nil {
			log.Error("OrderService.Delete Error when lock inventories: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}

		// Create inventory map for easy lookup
		inventoryMap := make(map[int]*entity.Inventory)
		for i := range lockedInventories {
			inv := &lockedInventories[i]
			inventoryMap[inv.ProductID] = inv
		}

		// Restore inventory quantities
		for _, item := range inventoryItems {
			inv := inventoryMap[item.ProductID]
			if inv == nil {
				log.Error("OrderService.Delete Error: inventory not found for productID ", item.ProductID)
				return error_utils.ErrorCode.DB_DOWN
			}

			quantityToRestore := item.Quantity
			newVersion := uuid.New().String()

			// Update inventory quantity
			err = s.inventoryRepo.UpdateQuantityWithVersionCommand(ctx, inv.ProductID, quantityToRestore, inv.Version, newVersion, tx)
			if err != nil {
				log.Error("OrderService.Delete Error when update inventory: " + err.Error())
				return error_utils.ErrorCode.DB_DOWN
			}

			// Create inventory history record for restoration
			inventoryHistory := &entity.InventoryHistory{
				ProductID:     inv.ProductID,
				Quantity:      quantityToRestore,
				FinalQuantity: inv.Quantity + quantityToRestore,
				ImporterName:  user.Username,
				ImportedAt:    time.Now(),
				Note:          "Hồi hàng về từ đơn xoá số " + strconv.Itoa(id),
			}
			err = s.inventoryHistoryRepo.CreateCommand(ctx, inventoryHistory, tx)
			if err != nil {
				log.Error("OrderService.Delete Error when create inventory history: " + err.Error())
				return error_utils.ErrorCode.DB_DOWN
			}

			// Update local inventory state
			inv.Quantity += quantityToRestore
			inv.Version = newVersion
		}

		// Delete the order
		err = s.orderRepo.DeleteByIDCommand(ctx, id, tx)
		if err != nil {
			log.Error("OrderService.Delete Error when delete order: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}

		// Commit transaction
		err = s.unitOfWork.Commit(tx)
		if err != nil {
			log.Error("OrderService.Delete Error when commit transaction: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}
	} else {
		// Delete the order
		err = s.orderRepo.DeleteByIDCommand(ctx, id, nil)
		if err != nil {
			log.Error("OrderService.Delete Error when delete order: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}
	}

	return ""
}
