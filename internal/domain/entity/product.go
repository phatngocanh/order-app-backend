package entity

type Product struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Spec string `db:"spec"`
	Type string `db:"type"`
}
