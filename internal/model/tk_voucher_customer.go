package model

import (
	"database/sql"
)

type TkVoucherCustomer struct {
	ID         int           `db:"id"`
	CustomerID sql.NullInt64 `db:"customer_id"`
	VoucherID  sql.NullInt64 `db:"voucher_id"`
	CreatedAt  sql.NullTime  `db:"created_at"`
	UpdatedAt  sql.NullTime  `db:"updated_at"`
	AmountLeft sql.NullInt64 `db:"amount_left"`
}
