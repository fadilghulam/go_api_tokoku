package model

import (
	"time"
)

const TableCustomerPointHistory = "tk.customer_point_history"

type CustomerPointHistory struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId    int64     `gorm:"column:customer_id" json:"customerId"`
	TransactionId int64     `gorm:"column:transaction_id;default:null" json:"transactionId"`
	ExchangeId    int64     `gorm:"column:exchange_id;default:null" json:"exchangeId"`
	DateTime      time.Time `gorm:"column:datetime;default:CURRENT_TIMESTAMP" json:"datetime"`
	Point         string    `gorm:"column:point" json:"point"`
	ExpiredDate   string    `gorm:"column:expired_date;default:null" json:"expired_date"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Type          string    `gorm:"column:type" json:"type"`
}

func (*CustomerPointHistory) TableName() string {
	return TableCustomerPointHistory
}
