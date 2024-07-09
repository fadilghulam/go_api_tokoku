package model

import (
	"database/sql"
)

type TkCartItem struct {
	ID             int           `db:"id"`
	CustomerID     sql.NullInt64 `db:"customer_id"`
	ItemExchangeID sql.NullInt64 `db:"item_exchange_id"`
	Point          sql.NullInt64 `db:"point"`
	CreatedAt      sql.NullTime  `db:"created_at"`
	UpdatedAt      sql.NullTime  `db:"updated_at"`
	Qty            sql.NullInt64 `db:"qty"`
}
