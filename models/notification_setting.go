package model

const TableNameTkNotificationSetting = "tk.notification_setting"

type TkNotificationSetting struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId    int64  `gorm:"column:customer_id;" json:"customerId"`
	OnTransaction *int16 `gorm:"column:on_transaction;default:1" json:"onTransaction"`
	OnNewPoint    *int16 `gorm:"column:on_new_point;default:1" json:"onNewPoint"`
	OnUsePoint    *int16 `gorm:"column:on_use_point;default:1" json:"onUsePoint"`
	Sound         *int16 `gorm:"column:sound;default:1" json:"sound"`
	Vibration     *int16 `gorm:"column:vibration;default:1" json:"vibration"`
	CreatedAt     string `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

type TkNotificationSettingInput struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId    int64  `gorm:"column:customer_id;" json:"customerId"`
	OnTransaction bool   `gorm:"column:on_transaction;" json:"onTransaction"`
	OnNewPoint    bool   `gorm:"column:on_new_point;" json:"onNewPoint"`
	OnUsePoint    bool   `gorm:"column:on_use_point;" json:"onUsePoint"`
	Sound         bool   `gorm:"column:sound;" json:"sound"`
	Vibration     bool   `gorm:"column:vibration;" json:"vibration"`
	CreatedAt     string `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*TkNotificationSetting) TableName() string {
	return TableNameTkNotificationSetting
}
