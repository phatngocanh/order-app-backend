package entity

import "time"

type OrderImage struct {
	ID        int       `db:"id"`
	OrderID   int       `db:"order_id"`
	ImageURL  string    `db:"image_url"`
	ImageType string    `db:"image_type"`
	CreatedAt time.Time `db:"created_at"`
}
