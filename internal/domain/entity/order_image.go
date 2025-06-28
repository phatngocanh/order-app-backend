package entity

type OrderImage struct {
	ID        int    `db:"id"`
	OrderID   int    `db:"order_id"`
	ImageURL  string `db:"image_url"`
	ImageType string `db:"image_type"`
	S3Key     string `db:"s3_key"`
}
