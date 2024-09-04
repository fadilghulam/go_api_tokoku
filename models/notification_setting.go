package model

import "time"

const TableNameTkNotificationSetting = "tk.notification_setting"

type TkNotificationSetting struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId    int64     `gorm:"column:customer_id;" json:"customerId"`
	OnTransaction int16     `gorm:"column:on_transaction;" json:"onTransaction"`
	OnNewPoint    int16     `gorm:"column:on_new_point;" json:"onNewPoint"`
	OnUsePoint    int16     `gorm:"column:on_use_point;" json:"onUsePoint"`
	Sound         int16     `gorm:"column:sound;" json:"sound"`
	Vibration     int16     `gorm:"column:vibration;" json:"vibration"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*TkNotificationSetting) TableName() string {
	return TableNameTkNotificationSetting
}
