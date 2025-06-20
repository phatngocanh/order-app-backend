package entity

import "time"

type Order struct {
	ID             int       `db:"id"`
	CustomerID     int       `db:"customer_id"`
	OrderDate      time.Time `db:"order_date"`
	DeliveryStatus string    `db:"delivery_status"`
	DebtStatus     string    `db:"debt_status"`
}
