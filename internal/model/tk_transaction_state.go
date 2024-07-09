package model

import (
	"database/sql"
)

type TkTransactionState struct {
	ID        int            `db:"id"`
	Name      sql.NullString `db:"name"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}
