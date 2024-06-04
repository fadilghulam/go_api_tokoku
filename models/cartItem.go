package model

import (
	"time"
)

const TableCartItem = "tk.cart_item"

type CartItem struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId     int64     `gorm:"column:customer_id" json:"customerId"`
	ItemExchangeId int64     `gorm:"column:item_exchange_id;default:null" json:"itemExchangeId"`
	Point          string    `gorm:"column:point" json:"point"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Qty            int16     `gorm:"column:qty" json:"qty"`
}

func (*CartItem) TableName() string {
	return TableCartItem
}
