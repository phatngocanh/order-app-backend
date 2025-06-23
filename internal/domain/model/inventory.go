package model

type UpdateInventoryQuantityRequest struct {
	Quantity int    `json:"quantity" binding:"required"`
	Note     string `json:"note"`
	Version  string `json:"version" binding:"required"`
}

type InventoryResponse struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Version   string `json:"version"`
}
