package entity

type Inventory struct {
	ID        int    `db:"id"`
	ProductID int    `db:"product_id"`
	Quantity  int    `db:"quantity"`
	Version   string `db:"version"`
}
