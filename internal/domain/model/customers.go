package model

type CreateCustomerRequest struct {
	Name    string `json:"name" binding:"required"`    // Tên khách hàng
	Phone   string `json:"phone" binding:"required"`   // Số điện thoại
	Address string `json:"address" binding:"required"` // Địa chỉ
}

type UpdateCustomerRequest struct {
	Name    string `json:"name"`    // Tên khách hàng
	Phone   string `json:"phone"`   // Số điện thoại
	Address string `json:"address"` // Địa chỉ
}

type CustomerResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`    // Tên khách hàng
	Phone   string `json:"phone"`   // Số điện thoại
	Address string `json:"address"` // Địa chỉ
}

type GetAllCustomersResponse struct {
	Customers []CustomerResponse `json:"customers"`
}

type GetOneCustomerResponse struct {
	Customer CustomerResponse `json:"customer"`
}
