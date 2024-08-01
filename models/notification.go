package model

import "time"

const TableNameTkNotification = "tk.notification"

type TkNotification struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId  int64     `gorm:"column:customer_id;" json:"customerId"`
	ProdukId    int32     `gorm:"column:produk_id" json:"produkId"`
	ItemId      int32     `gorm:"column:item_id" json:"itemId"`
	Title       string    `gorm:"column:title" json:"title"`
	Subtitle    string    `gorm:"column:subtitle" json:"subtitle"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*TkNotification) TableName() string {
	return TableNameTkNotification
}
