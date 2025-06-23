package model

type UpdateInventoryQuantityRequest struct {
	Quantity     int    `json:"quantity" binding:"required"`
	ImporterName string `json:"importer_name" binding:"required"`
	Note         string `json:"note"`
}

type InventoryResponse struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Version   string `json:"version"`
}
