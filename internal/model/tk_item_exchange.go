package model

import (
	"database/sql"
)

type TkItemExchange struct {
	ID            int            `db:"id"`
	ItemID        sql.NullInt64  `db:"item_id"`
	Point         sql.NullInt64  `db:"point"`
	DateStart     sql.NullTime   `db:"date_start"`
	DateEnd       sql.NullTime   `db:"date_end"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
	MaxExchange   sql.NullInt64  `db:"max_exchange"`
	About         sql.NullString `db:"about"`
	Detail        sql.NullString `db:"detail"`
	TermCondition sql.NullString `db:"term_condition"`
}
