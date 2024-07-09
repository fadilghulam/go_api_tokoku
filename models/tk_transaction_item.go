package model

import "time"

const TableNameTkTransactionItem = "tk.transaction_item"

type TkTransactionItem struct {
	ID                 int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TransactionStateID int64     `gorm:"column:transaction_state_id" json:"transaction_state_id"`
	CustomerID         int64     `gorm:"column:customer_id" json:"customer_id"`
	ReferenceID        int64     `gorm:"column:reference_id;default:null" json:"reference_id"`
	ReferenceName      string    `gorm:"column:reference_name;default:null" json:"reference_name"`
	TransactionDate    time.Time `gorm:"column:transaction_date" json:"transaction_date"`
	Provinsi           string    `gorm:"column:provinsi" json:"provinsi"`
	PronvisiID         int64     `gorm:"column:pronvisi_id;default:null" json:"pronvisi_id"`
	Kabupaten          string    `gorm:"column:kabupaten" json:"kabupaten"`
	KabupatenID        int64     `gorm:"column:kabupaten_id;default:null" json:"kabupaten_id"`
	Kecamatan          string    `gorm:"column:kecamatan" json:"kecamatan"`
	KecamatanID        int64     `gorm:"column:kecamatan_id;default:null" json:"kecamatan_id"`
	Kelurahan          string    `gorm:"column:kelurahan" json:"kelurahan"`
	KelurahanID        int64     `gorm:"column:kelurahan_id;default:null" json:"kelurahan_id"`
	SrID               int64     `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID            int64     `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	BranchID           int64     `gorm:"column:branch_id;default:null" json:"branch_id"`
	CreatedAt          time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	SyncKey            string    `gorm:"column:sync_key;default:now()" json:"sync_key"`
	Note               string    `gorm:"column:note" json:"note"`
}

func (*TkTransactionItem) TableName() string {
	return TableNameTkTransactionItem
}
