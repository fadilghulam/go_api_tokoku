package model

import (
	"database/sql"
)

type TkItemUnitTk struct {
	ID         int            `db:"id"`
	Name       sql.NullString `db:"name"`
	Multiplier sql.NullInt64  `db:"multiplier"`
	CreatedAt  sql.NullTime   `db:"created_at"`
	UpdatedAt  sql.NullTime   `db:"updated_at"`
}
