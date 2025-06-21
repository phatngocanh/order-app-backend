package model

type CreateProductRequest struct {
	Name          string `json:"name" binding:"required"`
	Spec          int    `json:"spec"`                              // Quy cách
	Type          string `json:"type"`                              // Loại hàng
	OriginalPrice int    `json:"original_price" binding:"required"` // Giá gốc (VND)
}

type UpdateProductRequest struct {
	ID            int    `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Spec          int    `json:"spec"`                              // Quy cách
	Type          string `json:"type"`                              // Loại hàng
	OriginalPrice int    `json:"original_price" binding:"required"` // Giá gốc (VND)
}

type ProductResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Spec          int    `json:"spec"`           // Quy cách
	Type          string `json:"type"`           // Loại hàng
	OriginalPrice int    `json:"original_price"` // Giá gốc (VND)
}

type GetAllProductsResponse struct {
	Products []ProductResponse `json:"products"`
}

type GetOneProductResponse struct {
	Product ProductResponse `json:"product"`
}
