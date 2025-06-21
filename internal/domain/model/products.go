package model

type CreateProductRequest struct {
	Name          string `json:"name" binding:"required"`           // Tên sản phẩm
	Spec          int    `json:"spec"`                              // Quy cách
	OriginalPrice int    `json:"original_price" binding:"required"` // Giá gốc của sản phẩm (VND)
}

type UpdateProductRequest struct {
	ID            int    `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`           // Tên sản phẩm
	Spec          int    `json:"spec"`                              // Quy cách
	OriginalPrice int    `json:"original_price" binding:"required"` // Giá gốc của sản phẩm (VND)
}

type ProductResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`           // Tên sản phẩm
	Spec          int    `json:"spec"`           // Quy cách
	OriginalPrice int    `json:"original_price"` // Giá gốc của sản phẩm (VND)
}

type GetAllProductsResponse struct {
	Products []ProductResponse `json:"products"`
}

type GetOneProductResponse struct {
	Product ProductResponse `json:"product"`
}
