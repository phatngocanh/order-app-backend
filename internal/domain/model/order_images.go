package model

type OrderImage struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"order_id"`
	ImageURL  string `json:"image_url"`
	ImageType string `json:"image_type"`
}

type UploadOrderImageResponse struct {
	OrderImage OrderImage `json:"orderImage"`
}

type GetOrderImagesResponse struct {
	OrderImages []OrderImage `json:"orderImages"`
}
