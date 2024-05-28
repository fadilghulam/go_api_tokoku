package model

import (
	"time"
)

const TableNameVoucherCustomer = "tk.voucher_customer"

type VoucherCustomer struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId int64     `gorm:"column:customer_id" json:"customerId"`
	VoucherId  int64     `gorm:"column:voucher_id" json:"voucherId"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*VoucherCustomer) TableName() string {
	return TableNameVoucherCustomer
}
