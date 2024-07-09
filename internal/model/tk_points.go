package model

import (
	"database/sql"
)

type TkPoints struct {
	ID        int            `db:"id"`
	BranchID  sql.NullInt64  `db:"branch_id"`
	ProdukID  sql.NullInt64  `db:"produk_id"`
	DateStart sql.NullTime   `db:"date_start"`
	DateEnd   sql.NullTime   `db:"date_end"`
	Value     sql.NullInt64  `db:"value"`
	Note      sql.NullString `db:"note"`
	CreatedBy sql.NullInt64  `db:"created_by"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}
