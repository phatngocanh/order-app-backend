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
	CHUA_GIAO       string
	DA_GIAO         string
	CHUA_THANH_TOAN string
	HOAN_THANH      string
}

var OrderDeliveryStatus = orderDeliveryStatus{
	CHUA_GIAO:       "CHUA_GIAO",
	DA_GIAO:         "DA_GIAO",
	CHUA_THANH_TOAN: "CHUA_THANH_TOAN",
	HOAN_THANH:      "HOAN_THANH",
}
