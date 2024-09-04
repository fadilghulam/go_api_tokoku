package model

import (
	"time"
)

const TableNameFCM = "tk.token_fcm_toko"

type TokenFcm struct {
	UserID            int32     `gorm:"column:user_id" json:"userId"`
	Token             string    `gorm:"column:token" json:"token"`
	AppName           string    `gorm:"column:app_name" json:"appName"`
	ApiKey            string    `gorm:"column:apiKey;default:null" json:"apiKey"`
	AppId             string    `gorm:"column:appId;default:null" json:"appId"`
	MessagingSenderId string    `gorm:"column:messagingSenderId;default:null" json:"messagingSenderId"`
	ProjectId         string    `gorm:"column:projectId;default:null" json:"projectId"`
	StorageBucket     string    `gorm:"column:storageBucket;default:null" json:"storageBucket"`
	IosClientId       string    `gorm:"column:iosClientId;default:null" json:"iosClientId"`
	IosBundleId       string    `gorm:"column:iosBundleId;default:null" json:"iosBundleId"`
	DeviceId          string    `gorm:"column:deviceId;default:null" json:"deviceId"`
	AppVersion        string    `gorm:"column:appVersion;default:null" json:"appVersion"`
	CreatedAt         time.Time `gorm:"column:created_at;not null;default:now()" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null;default:now()" json:"updatedAt"`
	CustomerId        int64     `gorm:"column:customer_id" json:"customerId"`
}

func (*TokenFcm) TableName() string {
	return TableNameFCM
}
