package model

import "time"

const TableNameTkNotification = "tk.notification"

type TkNotification struct {
	ID            string    `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerId    string    `gorm:"column:customer_id;" json:"customerId"`
	ProdukId      int32     `gorm:"column:produk_id;default:null" json:"produkId"`
	ItemId        int32     `gorm:"column:item_id;default:null" json:"itemId"`
	Title         string    `gorm:"column:title;default:null" json:"title"`
	Subtitle      string    `gorm:"column:subtitle;default:null" json:"subtitle"`
	Description   string    `gorm:"column:description;default:null" json:"description"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	ReferenceId   string    `gorm:"column:reference_id;default:null" json:"referenceId"`
	ReferenceName string    `gorm:"column:reference_name;default:null" json:"referenceName"`
	IsClose       int16     `gorm:"column:is_close;default:0" json:"isClose"`
}

func (*TkNotification) TableName() string {
	return TableNameTkNotification
}
