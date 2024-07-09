package model

import (
	"database/sql"
)

type TkDiscount struct {
	ID             int            `db:"id"`
	BranchID       sql.NullInt64  `db:"branch_id"`
	ProdukID       sql.NullInt64  `db:"produk_id"`
	DateStart      sql.NullTime   `db:"date_start"`
	DateEnd        sql.NullTime   `db:"date_end"`
	Nominal        sql.NullInt64  `db:"nominal"`
	CreatedBy      sql.NullInt64  `db:"created_by"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	CustomerTypeID sql.NullInt64  `db:"customer_type_id"`
	Note           sql.NullString `db:"note"`
}
