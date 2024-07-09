package model

import (
	"database/sql"
)

type TkCart struct {
	ID         int           `db:"id"`
	CustomerID sql.NullInt64 `db:"customer_id"`
	ProdukID   sql.NullInt64 `db:"produk_id"`
	Qty        sql.NullInt64 `db:"qty"`
	CreatedAt  sql.NullTime  `db:"created_at"`
	UpdatedAt  sql.NullTime  `db:"updated_at"`
	DateCart   sql.NullTime  `db:"date_cart"`
	StoreID    sql.NullInt64 `db:"store_id"`
	Harga      sql.NullInt64 `db:"harga"`
}
