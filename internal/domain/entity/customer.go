package entity

type Customer struct {
	ID      int    `db:"id"`
	Name    string `db:"name"`
	Phone   string `db:"phone"`
	Address string `db:"address"`
}

type customerLocationType struct {
	TINH      string
	THANH_PHO string
}

var CustomerLocationType = customerLocationType{
	TINH:      "TINH",
	THANH_PHO: "THANH_PHO",
}
