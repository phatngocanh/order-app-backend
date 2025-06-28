package service

import (
	"context"

	"github.com/pna/order-app-backend/internal/domain/model"
)

type StatisticsService interface {
	GetDashboardStats(ctx context.Context) (model.DashboardStatsResponse, string)
}
