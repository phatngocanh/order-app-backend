package entity

type Product struct {
	ID            int    `db:"id"`
	Name          string `db:"name"`           // Tên sản phẩm
	Spec          int    `db:"spec"`           // Quy cách
	OriginalPrice int    `db:"original_price"` // Giá gốc của sản phẩm (VND)
}
