package entity

import "time"

type InventoryHistory struct {
	ID           int       `db:"id"`
	ProductID    int       `db:"product_id"`
	Quantity     int       `db:"quantity"`
	ImporterName string    `db:"importer_name"`
	ImportedAt   time.Time `db:"imported_at"`
	Note         string    `db:"note"`
}
