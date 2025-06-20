package entity

type OrderItem struct {
	ID          int     `db:"id"`
	OrderID     int     `db:"order_id"`
	ProductID   int     `db:"product_id"`
	Quantity    int     `db:"quantity"`
	UnitPrice   float64 `db:"unit_price"`
	Discount    float64 `db:"discount"`
	FinalAmount float64 `db:"final_amount"`
}
