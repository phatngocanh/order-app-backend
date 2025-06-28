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

type InventoryWithProductResponse struct {
	ID        int         `json:"id"`
	ProductID int         `json:"product_id"`
	Quantity  int         `json:"quantity"`
	Version   string      `json:"version"`
	Product   ProductInfo `json:"product"`
}

type ProductInfo struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Spec          int    `json:"spec"`
	OriginalPrice int    `json:"original_price"`
}

type GetAllInventoryResponse struct {
	Inventories []InventoryWithProductResponse `json:"inventories"`
}
