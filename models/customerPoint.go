package model

import (
	"time"
)

const TableCustomerPoint = "tk.customer_point"

type CustomerPoint struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId int64     `gorm:"column:customer_id" json:"customerId"`
	Point      int64     `gorm:"column:point" json:"point"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*CustomerPoint) TableName() string {
	return TableCustomerPoint
}
