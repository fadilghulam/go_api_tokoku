package model

import "time"

const TableNameTkTransactionItemDetail = "tk.transaction_item_detail"

type TkTransactionItemDetail struct {
	ID                int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TransactionItemID int64     `db:"transaction_item_id" json:"transaction_item_id"`
	ItemID            int64     `db:"item_id" json:"item_id"`
	Qty               int64     `db:"qty" json:"qty"`
	SyncKey           string    `db:"sync_key" json:"sync_key"`
	CreatedAt         time.Time `db:"created_at;not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at;not null;default:now()" json:"updated_at"`
	Note              string    `db:"note" json:"note"`
	Point             int64     `db:"point" json:"point"`
}

func (*TkTransactionItemDetail) TableName() string {
	return TableNameTkTransactionItemDetail
}
