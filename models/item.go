package model

import (
	"time"
)

const TableNameItem = "tk.item"

type Item struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Image       string    `gorm:"column:image" json:"image"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*Item) TableName() string {
	return TableNameItem
}
