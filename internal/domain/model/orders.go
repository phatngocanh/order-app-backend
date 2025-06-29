package model

import "time"

type Order struct {
	ID                   int         `json:"id"`
	CustomerID           int         `json:"customer_id"`
	OrderDate            time.Time   `json:"order_date"`
	DeliveryStatus       string      `json:"delivery_status"`
	DebtStatus           string      `json:"debt_status"`
	StatusTransitionedAt *time.Time  `json:"status_transitioned_at"`
	OrderItems           []OrderItem `json:"order_items,omitempty"`
}

type OrderItem struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	OrderID       int    `json:"order_id"`
	ProductID     int    `json:"product_id"`
	NumberOfBoxes *int   `json:"number_of_boxes"`
	Spec          *int   `json:"spec"`
	Quantity      int    `json:"quantity"`
	SellingPrice  int    `json:"selling_price"`
	Discount      int    `json:"discount"`
	FinalAmount   *int   `json:"final_amount"`
	ExportFrom    string `json:"export_from"`
}

type CreateOrderRequest struct {
	CustomerID           int                `json:"customer_id" binding:"required"`      // Mã khách hàng
	OrderDate            time.Time          `json:"order_date" binding:"required"`       // Ngày đặt hàng
	DeliveryStatus       string             `json:"delivery_status" binding:"required"`  // Trạng thái giao hàng
	DebtStatus           *string            `json:"debt_status"`                         // Trạng thái công nợ
	StatusTransitionedAt *time.Time         `json:"status_transitioned_at"`              // Ngày chuyển trạng thái
	AdditionalCost       int                `json:"additional_cost"`                     // Chi phí phát sinh thêm (VND)
	AdditionalCostNote   *string            `json:"additional_cost_note"`                // Ghi chú cho chi phí phát sinh
	OrderItems           []OrderItemRequest `json:"order_items" binding:"required,dive"` // Danh sách sản phẩm trong đơn
}

type UpdateOrderRequest struct {
	ID                   int        `json:"id" binding:"required"`  // Mã đơn hàng
	CustomerID           int        `json:"customer_id"`            // Mã khách hàng
	OrderDate            time.Time  `json:"order_date"`             // Ngày đặt hàng
	DeliveryStatus       string     `json:"delivery_status"`        // Trạng thái giao hàng
	DebtStatus           *string    `json:"debt_status"`            // Trạng thái công nợ
	StatusTransitionedAt *time.Time `json:"status_transitioned_at"` // Ngày chuyển trạng thái
	AdditionalCost       *int       `json:"additional_cost"`        // Chi phí phát sinh thêm (VND)
	AdditionalCostNote   *string    `json:"additional_cost_note"`   // Ghi chú cho chi phí phát sinh
}

type OrderItemRequest struct {
	ProductID     int    `json:"product_id" binding:"required"`    // Mã sản phẩm
	NumberOfBoxes *int   `json:"number_of_boxes"`                  // Số thùng
	Spec          *int   `json:"spec"`                             // Quy cách mỗi thùng
	Quantity      int    `json:"quantity" binding:"required"`      // Số lượng cuối cùng
	SellingPrice  int    `json:"selling_price" binding:"required"` // Giá bán của sản phẩm (VND)
	Discount      int    `json:"discount"`                         // Chiết khấu (%)
	FinalAmount   *int   `json:"final_amount"`                     // Số tiền cuối cùng sau khi trừ chiết khấu (VND)
	Version       string `json:"version" binding:"required"`       // Version (UUID) của inventory để kiểm tra optimistic lock
	ExportFrom    string `json:"export_from" binding:"required"`   // Nguồn xuất: INVENTORY hoặc EXTERNAL
}

type OrderResponse struct {
	ID                   int                 `json:"id"`
	OrderDate            time.Time           `json:"order_date"`
	DeliveryStatus       string              `json:"delivery_status"`
	DebtStatus           *string             `json:"debt_status"`
	StatusTransitionedAt *time.Time          `json:"status_transitioned_at"`
	AdditionalCost       int                 `json:"additional_cost"`
	AdditionalCostNote   *string             `json:"additional_cost_note"`
	Customer             CustomerResponse    `json:"customer"`
	OrderItems           []OrderItemResponse `json:"order_items,omitempty"`
	Images               []OrderImage        `json:"images,omitempty"`
	TotalAmount          *int                `json:"total_amount,omitempty"`
	ProductCount         *int                `json:"product_count,omitempty"`
	// Profit/Loss fields for total order
	TotalProfitLoss           *int     `json:"total_profit_loss,omitempty"`            // Total profit/loss for the order
	TotalProfitLossPercentage *float64 `json:"total_profit_loss_percentage,omitempty"` // Total profit/loss percentage for the order
}

type OrderItemResponse struct {
	ID            int    `json:"id"`
	ProductName   string `json:"product_name"`
	OrderID       int    `json:"order_id"`
	ProductID     int    `json:"product_id"`
	NumberOfBoxes *int   `json:"number_of_boxes"`
	Spec          *int   `json:"spec"`
	Quantity      int    `json:"quantity"`
	SellingPrice  int    `json:"selling_price"`
	Discount      int    `json:"discount"`
	FinalAmount   *int   `json:"final_amount"`
	ExportFrom    string `json:"export_from"`
	// Profit/Loss fields
	OriginalPrice        *int     `json:"original_price,omitempty"`         // Product's original price
	ProfitLoss           *int     `json:"profit_loss,omitempty"`            // Profit/Loss amount for this item
	ProfitLossPercentage *float64 `json:"profit_loss_percentage,omitempty"` // Profit/Loss percentage for this item
}

type GetAllOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
}

type GetOneOrderResponse struct {
	Order OrderResponse `json:"order"`
}
