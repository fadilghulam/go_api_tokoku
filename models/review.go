package model

import (
	"time"
)

const TableNameReview = "tk.review"

type Review struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId  int64     `gorm:"column:customer_id" json:"customerId"`
	OrderId     int64     `gorm:"column:order_id;default:null" json:"orderId"`
	OrderItemId int64     `gorm:"column:order_item_id;default:null" json:"orderItemId"`
	Rating      int16     `gorm:"column:rating" json:"rating"`
	Description string    `gorm:"column:description;default:null" json:"description"`
	Photo       string    `gorm:"column:photo;default:null" json:"photo"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*Review) TableName() string {
	return TableNameReview
}
