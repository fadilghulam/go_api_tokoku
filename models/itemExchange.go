package model

import (
	"time"
)

const TableNameItemExchange = "tk.item_exchange"

type ItemExchange struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ItemId    int64     `gorm:"column:item_id" json:"item_id"`
	Point     int16     `gorm:"column:point" json:"point"`
	DateStart string    `gorm:"column:date_start;not null;default:now()" json:"date_start"`
	DateEnd   string    `gorm:"column:date_end;not null;default:now()" json:"date_end"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (*ItemExchange) TableName() string {
	return TableNameItemExchange
}
