package model

import (
	"database/sql"
)

type TkUnitMapping struct {
	ID         int           `db:"id"`
	ProdukID   sql.NullInt64 `db:"produk_id"`
	ItemID     sql.NullInt64 `db:"item_id"`
	ItemUnitID sql.NullInt64 `db:"item_unit_id"`
	CreatedAt  sql.NullTime  `db:"created_at"`
	UpdatedAt  sql.NullTime  `db:"updated_at"`
}
