package entity

import "time"

type Order struct {
	ID                   int        `db:"id"`
	CustomerID           int        `db:"customer_id"`
	OrderDate            time.Time  `db:"order_date"`
	DeliveryStatus       string     `db:"delivery_status"`
	DebtStatus           *string    `db:"debt_status"`
	StatusTransitionedAt *time.Time `db:"status_transitioned_at"`
	TotalOriginalCost    int        `db:"total_original_cost"`
	TotalSalesRevenue    int        `db:"total_sales_revenue"`
	AdditionalCost       int        `db:"additional_cost"`
	AdditionalCostNote   *string    `db:"additonal_cost_note"`
}

type orderDeliveryStatus struct {
	PENDING   string
	DELIVERED string
	UNPAID    string
	COMPLETED string
}

var OrderDeliveryStatus = orderDeliveryStatus{
	PENDING:   "PENDING",
	DELIVERED: "DELIVERED",
	UNPAID:    "UNPAID",
	COMPLETED: "COMPLETED",
}
