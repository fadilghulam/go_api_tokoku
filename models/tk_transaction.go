package model

import "time"

const TableNameTkTransaction = "tk.transaction"

type TkTransaction struct {
	ID                 int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TransactionStateID int64     `gorm:"column:transaction_state_id" json:"transaction_state_id"`
	CustomerID         int64     `gorm:"column:customer_id" json:"customer_id"`
	ReferenceID        int64     `gorm:"column:reference_id;default:null" json:"reference_id"`
	ReferenceName      string    `gorm:"column:reference_name;default:null" json:"reference_name"`
	TransactionDate    time.Time `gorm:"column:transaction_date" json:"transaction_date"`
	TotalTransaction   int64     `gorm:"column:total_transaction" json:"total_transaction"`
	Provinsi           string    `gorm:"column:provinsi;default:null" json:"provinsi"`
	PronvisiID         int64     `gorm:"column:pronvisi_id;default:null" json:"pronvisi_id"`
	Kabupaten          string    `gorm:"column:kabupaten;default:null" json:"kabupaten"`
	KabupatenID        int64     `gorm:"column:kabupaten_id;default:null" json:"kabupaten_id"`
	Kecamatan          string    `gorm:"column:kecamatan;default:null" json:"kecamatan"`
	KecamatanID        int64     `gorm:"column:kecamatan_id;default:null" json:"kecamatan_id"`
	Kelurahan          string    `gorm:"column:kelurahan;default:null" json:"kelurahan"`
	KelurahanID        int64     `gorm:"column:kelurahan_id;default:null" json:"kelurahan_id"`
	SrID               int64     `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID            int64     `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	BranchID           int64     `gorm:"column:branch_id;default:null" json:"branch_id"`
	CreatedAt          time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	SyncKey            string    `gorm:"column:sync_key;default:now()" json:"sync_key"`
	Note               string    `gorm:"column:note;default:null" json:"note"`
	StoreID            int64     `gorm:"column:store_id;default:null" json:"store_id"`
	EstimateDate       string    `gorm:"column:estimate_date;default:null" json:"estimate_date"`
}

func (*TkTransaction) TableName() string {
	return TableNameTkTransaction
}
