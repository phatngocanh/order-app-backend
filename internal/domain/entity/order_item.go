package entity

type OrderItem struct {
	ID            int    `db:"id"`
	OrderID       int    `db:"order_id"`
	ProductID     int    `db:"product_id"`
	NumberOfBoxes *int   `db:"number_of_boxes"` // Số thùng
	Spec          *int   `db:"spec"`            // Quy cách mỗi thùng
	Quantity      int    `db:"quantity"`        // Số lượng cuối cùng
	SellingPrice  int    `db:"selling_price"`   // Giá bán của sản phẩm (VND)
	OriginalPrice int    `db:"original_price"`  // Giá gốc của sản phẩm (VND)
	Discount      int    `db:"discount"`        // Chiết khấu (%)
	FinalAmount   *int   `db:"final_amount"`    // Số tiền cuối cùng sau khi trừ chiết khấu (VND)
	ExportFrom    string `db:"export_from"`     // Nguồn xuất: từ xưởng hoặc từ bên ngoài
}

type orderExportFrom struct {
	INVENTORY string
	EXTERNAL  string
}

var OrderExportFrom = orderExportFrom{
	INVENTORY: "INVENTORY",
	EXTERNAL:  "EXTERNAL",
}
