package model

import (
	"database/sql"
)

type TkVoucher struct {
	ID           int             `db:"id"`
	ProdukID     sql.NullString  `db:"produk_id"`
	DateStart    sql.NullTime    `db:"date_start"`
	DateEnd      sql.NullTime    `db:"date_end"`
	Code         sql.NullString  `db:"code"`
	Icon         sql.NullString  `db:"icon"`
	Diskon       sql.NullFloat64 `db:"diskon"`
	IsPercentage sql.NullInt64   `db:"is_percentage"`
	MinCost      sql.NullFloat64 `db:"min_cost"`
	MaxDiskon    sql.NullFloat64 `db:"max_diskon"`
	CreatedAt    sql.NullTime    `db:"created_at"`
	UpdatedAt    sql.NullTime    `db:"updated_at"`
	Note         sql.NullString  `db:"note"`
	Amount       sql.NullInt64   `db:"amount"`
}
