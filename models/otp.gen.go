package model

import (
	"time"
)

const TableNameOtp = "otp"

// Otp mapped from table <otp>
type Otp struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AppName     string    `gorm:"column:app_name" json:"appName"`
	Description string    `gorm:"column:description;default:null" json:"description"`
	Type        string    `gorm:"column:type" json:"type"`
	SendTo      string    `gorm:"column:send_to;not null" json:"sendTo"`
	Otp         string    `gorm:"column:otp;not null" json:"otp"`
	UserID      int32     `gorm:"column:user_id" json:"userId"`
	ExpiredAt   time.Time `gorm:"column:expired_at" json:"expiredAt"`
	ConfirmedAt string    `gorm:"column:confirmed_at;default:null" json:"confirmedAt"`
	CreatedAt   time.Time `gorm:"column:created_at;default:now()" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:now()" json:"updatedAt"`
	DeletedAt   time.Time `gorm:"column:deleted_at;default:null" json:"deletedAt"`
	Label       string    `gorm:"column:label;not null" json:"label"`
}

// TableName Otp's table name
func (*Otp) TableName() string {
	return TableNameOtp
}
