package serviceimplement

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pna/order-app-backend/internal/bean"
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
	orderImageRepo       repository.OrderImageRepository
	s3Service            bean.S3Service
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	inventoryRepo repository.InventoryRepository,
	inventoryHistoryRepo repository.InventoryHistoryRepository,
	userRepo repository.UserRepository,
	unitOfWork repository.UnitOfWork,
	productRepo repository.ProductRepository,
	customerRepo repository.CustomerRepository,
	orderImageRepo repository.OrderImageRepository,
	s3Service bean.S3Service) service.OrderService {
	return &OrderService{
		orderRepo:            orderRepo,
		orderItemRepo:        orderItemRepo,
		inventoryRepo:        inventoryRepo,
		inventoryHistoryRepo: inventoryHistoryRepo,
		userRepo:             userRepo,
		unitOfWork:           unitOfWork,
		productRepo:          productRepo,
		customerRepo:         customerRepo,
		orderImageRepo:       orderImageRepo,
		s3Service:            s3Service,
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

// Helper to calculate total original cost and total sales revenue for order items
func (s *OrderService) calculateOrderCostAndRevenue(ctx context.Context, orderItems []model.OrderItemRequest) (totalOriginalCost int, totalSalesRevenue int, err string) {
	totalOriginalCost = 0
	totalSalesRevenue = 0

	for _, item := range orderItems {
		// Get product to get original price
		product, err := s.productRepo.GetOneByIDQuery(ctx, item.ProductID, nil)
		if err != nil {
			log.Error("OrderService.calculateOrderCostAndRevenue Error fetching product: " + err.Error())
			return 0, 0, error_utils.ErrorCode.DB_DOWN
		}
		if product == nil {
			log.Error("OrderService.calculateOrderCostAndRevenue Error: product not found for ID: ", item.ProductID)
			return 0, 0, error_utils.ErrorCode.NOT_FOUND
		}

		// Calculate original cost
		originalCost := item.Quantity * product.OriginalPrice
		totalOriginalCost += originalCost

		// Calculate final revenue (after discount)
		sellingRevenue := item.Quantity * item.SellingPrice
		discountAmount := (sellingRevenue * item.Discount) / 100
		finalRevenue := sellingRevenue - discountAmount
		totalSalesRevenue += finalRevenue
	}

	return totalOriginalCost, totalSalesRevenue, ""
}

func (s *OrderService) GetAll(ctx context.Context, userID int, customerID int, deliveryStatuses string, sortBy string) (model.GetAllOrdersResponse, string) {
	orders, err := s.orderRepo.GetAllWithFiltersQuery(ctx, customerID, deliveryStatuses, sortBy, nil)
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
		totalAmount += o.AdditionalCost
		// Calculate profit/loss from stored cost and revenue values
		totalProfitLoss := o.TotalSalesRevenue - o.TotalOriginalCost + o.AdditionalCost
		totalProfitLossPercentage := 0.0
		if o.TotalOriginalCost > 0 {
			totalProfitLossPercentage = float64(totalProfitLoss) / float64(o.TotalOriginalCost) * 100
		}

		resp.Orders = append(resp.Orders, model.OrderResponse{
			ID:                   o.ID,
			OrderDate:            o.OrderDate,
			DeliveryStatus:       o.DeliveryStatus,
			DebtStatus:           o.DebtStatus,
			StatusTransitionedAt: o.StatusTransitionedAt,
			AdditionalCost:       o.AdditionalCost,
			AdditionalCostNote:   o.AdditionalCostNote,
			Customer: model.CustomerResponse{
				ID:      customer.ID,
				Name:    customer.Name,
				Phone:   customer.Phone,
				Address: customer.Address,
			},
			OrderItems:                nil, // Omit order items in GetAll
			TotalAmount:               &totalAmount,
			ProductCount:              &productCount,
			TotalProfitLoss:           &totalProfitLoss,
			TotalProfitLossPercentage: &totalProfitLossPercentage,
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
	totalOriginalCost := 0
	totalProfitLoss := 0

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

		// Calculate profit/loss for this item
		originalCost := item.Quantity * item.OriginalPrice
		sellingRevenue := item.Quantity * item.SellingPrice
		discountAmount := (sellingRevenue * item.Discount) / 100
		finalRevenue := sellingRevenue - discountAmount
		profitLoss := finalRevenue - originalCost
		profitLossPercentage := 0.0
		if originalCost > 0 {
			profitLossPercentage = float64(profitLoss) / float64(originalCost) * 100
		}

		// Accumulate totals
		totalOriginalCost += originalCost
		totalProfitLoss += profitLoss

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
			// Profit/Loss fields
			OriginalPrice:        &item.OriginalPrice,
			ProfitLoss:           &profitLoss,
			ProfitLossPercentage: &profitLossPercentage,
		})
	}

	totalAmount, productCount := calculateOrderAmountsAndProductCount(orderItems)
	totalAmount += order.AdditionalCost

	// Use stored values for total order profit/loss
	totalProfitLoss = order.TotalSalesRevenue - order.TotalOriginalCost + order.AdditionalCost
	totalProfitLossPercentage := 0.0
	if order.TotalOriginalCost > 0 {
		totalProfitLossPercentage = float64(totalProfitLoss) / float64(order.TotalOriginalCost) * 100
	}

	// Fetch order images and generate signed URLs
	orderImages, err := s.orderImageRepo.GetAllByOrderIDQuery(ctx, order.ID, nil)
	if err != nil {
		log.Error("OrderService.GetOne Error fetching order images: " + err.Error())
		// Continue without images rather than failing the entire request
		orderImages = make([]entity.OrderImage, 0)
	}

	// Convert images to response model with signed URLs
	var imageResponses []model.OrderImage
	if len(orderImages) > 0 {
		for _, img := range orderImages {
			// Generate a fresh signed URL for each image
			signedURL, err := s.s3Service.GenerateSignedDownloadURL(ctx, img.S3Key, 20*time.Second)
			if err != nil {
				log.Error("OrderService.GetOne Error generating signed URL for image: " + err.Error())
				// Continue with other images even if one fails
				signedURL = ""
			}

			imageResponses = append(imageResponses, model.OrderImage{
				ID:       img.ID,
				OrderID:  img.OrderID,
				ImageURL: signedURL,
				S3Key:    img.S3Key,
			})
		}
	}

	resp := model.GetOneOrderResponse{Order: model.OrderResponse{
		ID:                   order.ID,
		OrderDate:            order.OrderDate,
		DeliveryStatus:       order.DeliveryStatus,
		DebtStatus:           order.DebtStatus,
		StatusTransitionedAt: order.StatusTransitionedAt,
		AdditionalCost:       order.AdditionalCost,
		AdditionalCostNote:   order.AdditionalCostNote,
		Customer: model.CustomerResponse{
			ID:      customer.ID,
			Name:    customer.Name,
			Phone:   customer.Phone,
			Address: customer.Address,
		},
		OrderItems:   orderItemResponses,
		Images:       imageResponses,
		TotalAmount:  &totalAmount,
		ProductCount: &productCount,
		// Profit/Loss fields for total order
		TotalProfitLoss:           &totalProfitLoss,
		TotalProfitLossPercentage: &totalProfitLossPercentage,
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

	// Calculate total original cost and total sales revenue
	totalOriginalCost, totalSalesRevenue, errCode := s.calculateOrderCostAndRevenue(ctx, req.OrderItems)
	if errCode != "" {
		return errCode
	}

	orderEntity := entity.Order{
		CustomerID:         req.CustomerID,
		OrderDate:          req.OrderDate,
		DeliveryStatus:     req.DeliveryStatus,
		DebtStatus:         req.DebtStatus,
		TotalOriginalCost:  totalOriginalCost,
		TotalSalesRevenue:  totalSalesRevenue,
		AdditionalCost:     req.AdditionalCost,
		AdditionalCostNote: req.AdditionalCostNote,
	}
	now := time.Now()
	orderEntity.StatusTransitionedAt = &now

	err = s.orderRepo.CreateCommand(ctx, &orderEntity, tx)
	if err != nil {
		log.Error("OrderService.Create Error when create order: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// Validate that each product has at most 2 order items (1 from inventory, 1 from external)
	productOrderItemCount := make(map[int]map[string]int) // productID -> exportFrom -> count
	for _, item := range req.OrderItems {
		if productOrderItemCount[item.ProductID] == nil {
			productOrderItemCount[item.ProductID] = make(map[string]int)
		}
		productOrderItemCount[item.ProductID][item.ExportFrom]++

		if productOrderItemCount[item.ProductID][item.ExportFrom] > 1 {
			log.Error("OrderService.Create Error: product ", item.ProductID, " has more than 1 order item from ", item.ExportFrom)
			return error_utils.ErrorCode.DUPLICATE_ORDER_ITEMS
		}
	}

	for _, item := range req.OrderItems {
		inv := inventoryMap[item.ProductID]
		quantityToExport := item.Quantity

		// Only check version for inventory items
		if item.ExportFrom == entity.OrderExportFrom.INVENTORY {
			itemVersion := item.Version
			if inv.Version != itemVersion {
				log.Error("OrderService.Create Error: inventory version mismatch for productID ", item.ProductID)
				return error_utils.ErrorCode.INVENTORY_VERSION_MISMATCH
			}
		}

		product, err := s.productRepo.GetOneByIDQuery(ctx, item.ProductID, nil)
		if err != nil {
			log.Error("OrderService.Create Error fetching product: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}

		// Calculate final amount for the item
		itemTotal := quantityToExport * item.SellingPrice
		discountAmount := (itemTotal * item.Discount) / 100
		finalAmount := itemTotal - discountAmount

		itemEntity := entity.OrderItem{
			ProductID:     item.ProductID,
			NumberOfBoxes: item.NumberOfBoxes,
			Spec:          item.Spec,
			Quantity:      quantityToExport,
			SellingPrice:  item.SellingPrice,
			OriginalPrice: product.OriginalPrice,
			Discount:      item.Discount,
			FinalAmount:   &finalAmount,
			OrderID:       orderEntity.ID,
			ExportFrom:    item.ExportFrom,
		}

		// Handle based on export source
		if item.ExportFrom == entity.OrderExportFrom.INVENTORY {
			// Check if inventory has enough quantity
			if inv.Quantity < quantityToExport {
				log.Error("OrderService.Create Error: inventory quantity exceeded for productID ", item.ProductID)
				return error_utils.ErrorCode.INVENTORY_QUANTITY_EXCEEDED
			}

			// Update inventory quantity
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

			// Update local inventory state
			inv.Quantity -= quantityToExport
			inv.Version = newVersion
		} else if item.ExportFrom != entity.OrderExportFrom.EXTERNAL {
			// Invalid export source
			log.Error("OrderService.Create Error: invalid export_from value for productID ", item.ProductID)
			return error_utils.ErrorCode.BAD_REQUEST
		}
		// For EXTERNAL source, no inventory operations are needed - items will be sourced from external suppliers

		// Create order item
		err = s.orderItemRepo.CreateCommand(ctx, &itemEntity, tx)
		if err != nil {
			log.Error("OrderService.Create Error when create order item: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
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
	if req.DebtStatus != nil {
		existing.DebtStatus = req.DebtStatus
	}
	if req.AdditionalCost != nil {
		existing.AdditionalCost = *req.AdditionalCost
	}
	if req.AdditionalCostNote != nil {
		existing.AdditionalCostNote = req.AdditionalCostNote
	}

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
