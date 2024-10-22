package model

import "time"

const TableNameTkNotification = "tk.notification"

type TkNotification struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId    int64     `gorm:"column:customer_id;" json:"customerId"`
	ProdukId      int32     `gorm:"column:produk_id;default:null" json:"produkId"`
	ItemId        int32     `gorm:"column:item_id;default:null" json:"itemId"`
	Title         string    `gorm:"column:title;default:null" json:"title"`
	Subtitle      string    `gorm:"column:subtitle;default:null" json:"subtitle"`
	Description   string    `gorm:"column:description;default:null" json:"description"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	ReferenceId   int64     `gorm:"column:reference_id;default:null" json:"reference_id"`
	ReferenceName string    `gorm:"column:reference_name;default:null" json:"reference_name"`
}

func (*TkNotification) TableName() string {
	return TableNameTkNotification
}
