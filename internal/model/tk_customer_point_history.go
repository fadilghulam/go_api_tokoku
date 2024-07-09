package model

import (
	"database/sql"
)

type TkCustomerPointHistory struct {
	ID            int            `db:"id"`
	CustomerID    sql.NullInt64  `db:"customer_id"`
	TransactionID sql.NullInt64  `db:"transaction_id"`
	ExchangeID    sql.NullInt64  `db:"exchange_id"`
	Datetime      sql.NullTime   `db:"datetime"`
	Point         sql.NullInt64  `db:"point"`
	ExpiredDate   sql.NullTime   `db:"expired_date"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
	Type          sql.NullString `db:"type"`
}
