package entity

import "time"

type Order struct {
	ID             int       `db:"id"`
	CustomerID     int       `db:"customer_id"`
	OrderDate      time.Time `db:"order_date"`
	DeliveryStatus string    `db:"delivery_status"`
	DebtStatus     string    `db:"debt_status"`
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
