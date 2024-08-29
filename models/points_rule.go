package model

import (
	"time"
)

const TableNamePointsRule = "tk.points_rule"

type PointsRule struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	BranchId  int32     `gorm:"column:branch_id" json:"branch_id"`
	ProdukId  int32     `gorm:"column:produk_id" json:"produk_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Type      string    `gorm:"column:type" json:"type"`
	SubType   string    `gorm:"column:sub_type" json:"sub_type"`
	Min       int16     `gorm:"column:min" json:"min"`
	Point     int16     `gorm:"column:point" json:"point"`
	MaxPoint  int16     `gorm:"column:max_point" json:"max_point"`
	Kelipatan int16     `gorm:"column:kelipatan" json:"kelipatan"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	DateStart time.Time `gorm:"column:date_start;not null;default:now()" json:"date_start"`
	DateEnd   time.Time `gorm:"column:date_end;not null;default:now()" json:"date_end"`
}

func (*PointsRule) TableName() string {
	return TableNamePointsRule
}
