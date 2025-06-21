package entity

type Product struct {
	ID            int    `db:"id"`
	Name          string `db:"name"`
	Spec          int    `db:"spec"`           // Quy cách
	Type          string `db:"type"`           // Loại hàng
	OriginalPrice int    `db:"original_price"` // Giá gốc của sản phẩm (VND)
}
