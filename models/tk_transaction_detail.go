package model

import "time"

const TableNameTkTransactionDetail = "tk.transaction_detail"

type TkTransactionDetail struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TransactionID int64     `gorm:"column:transaction_id" json:"transaction_id"`
	ProdukID      int64     `gorm:"column:produk_id" json:"produk_id"`
	Qty           int64     `gorm:"column:qty" json:"qty"`
	Harga         float64   `gorm:"column:harga" json:"harga"`
	Diskon        float64   `gorm:"column:diskon" json:"diskon"`
	Condition     string    `gorm:"column:condition" json:"condition"`
	Pita          string    `gorm:"column:pita" json:"pita"`
	SyncKey       string    `gorm:"column:sync_key" json:"sync_key"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Note          string    `gorm:"column:note" json:"note"`
	Point         int64     `gorm:"column:point" json:"point"`
}

func (*TkTransactionDetail) TableName() string {
	return TableNameTkTransactionDetail
}
