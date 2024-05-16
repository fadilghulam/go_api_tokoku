package model

import (
	"time"
)

const TableNameCart = "tk.cart"

type Cart struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId int64     `gorm:"column:customer_id" json:"customer_id"`
	ProdukId   int32     `gorm:"column:produk_id" json:"produk_id"`
	Qty        int64     `gorm:"column:qty" json:"qty"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	DateCart   time.Time `gorm:"column:date_cart;not null;default:current_date" json:"date_cart"`
}

func (*Cart) TableName() string {
	return TableNameCart
}
