package model

import (
	"time"
)

const TableNameAdvertisement = "tk.advertisement"

type Advertisement struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ProdukId    int32     `gorm:"column:produk_id" json:"produk_id"`
	ItemId      int32     `gorm:"column:item_id" json:"item_id"`
	Title       string    `gorm:"column:title" json:"title"`
	Description string    `gorm:"column:description" json:"description"`
	Image       string    `gorm:"column:image" json:"image"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	DateStart   time.Time `gorm:"column:date_start;not null;default:now()" json:"date_start"`
	DateEnd     time.Time `gorm:"column:date_end;not null;default:now()" json:"date_end"`
}

func (*Advertisement) TableName() string {
	return TableNameAdvertisement
}
