package model

import "time"

type CreateInventoryHistoryRequest struct {
	ProductID    int    `json:"product_id" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required"`
	ImporterName string `json:"importer_name" binding:"required"`
	Note         string `json:"note"`
}

type InventoryHistoryResponse struct {
	ID            int       `json:"id"`
	ProductID     int       `json:"product_id"`
	Quantity      int       `json:"quantity"`
	FinalQuantity int       `json:"final_quantity"`
	ImporterName  string    `json:"importer_name"`
	ImportedAt    time.Time `json:"imported_at"`
	Note          string    `json:"note"`
	ReferenceID   *int      `json:"reference_id,omitempty"`
}

type GetAllInventoryHistoriesResponse struct {
	InventoryHistories []InventoryHistoryResponse `json:"inventory_histories"`
}
