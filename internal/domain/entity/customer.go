package entity

type Customer struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	Phone        string `db:"phone"`
	Address      string `db:"address"`
	LocationType string `db:"location_type"`
}
