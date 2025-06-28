package model

type OrderImage struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"order_id"`
	ImageURL  string `json:"image_url"`
	ImageType string `json:"image_type"`
	S3Key     string `json:"s3_key"`
}

type UploadOrderImageResponse struct {
	OrderImage OrderImage `json:"orderImage"`
}

type GetOrderImagesResponse struct {
	OrderImages []OrderImage `json:"orderImages"`
}

type GenerateSignedUploadURLResponse struct {
	SignedURL string `json:"signed_url"`
	S3Key     string `json:"s3_key"`
}
