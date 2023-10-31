package model

type Account struct {
	ID            string  `db:"id"`
	OwnerID       string  `db:"owner_id"`
	AccountType   int64   `db:"account_type"`
	Balance       float64 `db:"balance"`
	AccountNo     []byte  `db:"account_no"`
	AccountNoHash []byte  `db:"account_no_hash"`
	PIN           string  `db:"pin"`
	RowHash       []byte  `db:"row_hash"`
}
