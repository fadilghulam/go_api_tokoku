package model

import "time"

const TableNameTkItem = "tk.item"

type TkItem struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Image       string    `gorm:"column:image" json:"image"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*TkItem) TableName() string {
	return TableNameTkItem
}
