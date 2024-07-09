package model

import (
	"database/sql"
)

type TkKritikSaran struct {
	ID         int            `db:"id"`
	CustomerID sql.NullInt64  `db:"customer_id"`
	Fitur      sql.NullString `db:"fitur"`
	Note       sql.NullString `db:"note"`
	Lampiran   sql.NullString `db:"lampiran"`
	CreatedAt  sql.NullTime   `db:"created_at"`
	UpdatedAt  sql.NullTime   `db:"updated_at"`
}
