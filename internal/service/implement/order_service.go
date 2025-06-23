package serviceimplement

import (
	"context"
	"database/sql"
	"time"

	"github.com/pna/order-app-backend/internal/domain/entity"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/repository"
	"github.com/pna/order-app-backend/internal/service"
)

type OrderService struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	unitOfWork    repository.UnitOfWork
}

func NewOrderService(orderRepo repository.OrderRepository, orderItemRepo repository.OrderItemRepository, unitOfWork repository.UnitOfWork) service.OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		unitOfWork:    unitOfWork,
	}
}

func (s *OrderService) GetAll(ctx context.Context) (model.GetAllOrdersResponse, error) {
	orders, err := s.orderRepo.GetAllQuery(ctx, nil)
	if err != nil {
		return model.GetAllOrdersResponse{}, err
	}
	resp := model.GetAllOrdersResponse{Orders: make([]model.OrderResponse, 0, len(orders))}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, model.OrderResponse{
			ID:                   o.ID,
			CustomerID:           o.CustomerID,
			OrderDate:            o.OrderDate,
			DeliveryStatus:       o.DeliveryStatus,
			DebtStatus:           o.DebtStatus,
			StatusTransitionedAt: o.StatusTransitionedAt,
			OrderItems:           []model.OrderItemResponse{}, // Not loaded here
		})
	}
	return resp, nil
}

func (s *OrderService) GetOne(ctx context.Context, id int) (model.GetOneOrderResponse, error) {
	order, err := s.orderRepo.GetOneByIDQuery(ctx, id, nil)
	if err != nil {
		return model.GetOneOrderResponse{}, err
	}
	if order == nil {
		return model.GetOneOrderResponse{}, nil
	}
	resp := model.GetOneOrderResponse{Order: model.OrderResponse{
		ID:                   order.ID,
		CustomerID:           order.CustomerID,
		OrderDate:            order.OrderDate,
		DeliveryStatus:       order.DeliveryStatus,
		DebtStatus:           order.DebtStatus,
		StatusTransitionedAt: order.StatusTransitionedAt,
		OrderItems:           []model.OrderItemResponse{}, // Not loaded here
	}}
	return resp, nil
}

func (s *OrderService) Create(ctx context.Context, req model.CreateOrderRequest) error {
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			s.unitOfWork.Rollback(tx)
		}
	}()

	orderEntity := entity.Order{
		CustomerID:     req.CustomerID,
		OrderDate:      req.OrderDate,
		DeliveryStatus: req.DeliveryStatus,
		DebtStatus:     req.DebtStatus,
	}
	now := time.Now()
	orderEntity.StatusTransitionedAt = &now

	err = s.orderRepo.CreateCommand(ctx, &orderEntity, tx)
	if err != nil {
		return err
	}

	for _, item := range req.OrderItems {
		itemEntity := entity.OrderItem{
			ProductID:     item.ProductID,
			NumberOfBoxes: item.NumberOfBoxes,
			Spec:          item.Spec,
			Quantity:      item.Quantity,
			SellingPrice:  item.SellingPrice,
			Discount:      item.Discount,
			FinalAmount:   item.FinalAmount,
			OrderID:       orderEntity.ID,
		}
		err = s.orderItemRepo.CreateCommand(ctx, &itemEntity, tx)
		if err != nil {
			return err
		}
	}

	return s.unitOfWork.Commit(tx)
}

func (s *OrderService) Update(ctx context.Context, req model.UpdateOrderRequest) error {
	existing, err := s.orderRepo.GetOneByIDQuery(ctx, req.ID, nil)
	if err != nil {
		return err
	}
	if existing == nil {
		return sql.ErrNoRows
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

	return s.orderRepo.UpdateCommand(ctx, existing, nil)
}
