package model

import (
	"time"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID              int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	FullName        string    `gorm:"column:full_name;not null" json:"full_name"`
	Username        string    `gorm:"column:username;not null" json:"username"`
	Password        string    `gorm:"column:password;not null" json:"password"`
	Images          string    `gorm:"column:images" json:"images"`
	LevelIDOld      int32     `gorm:"column:level_id_old" json:"level_id_old"`
	IsAktif         string    `gorm:"column:is_aktif;not null" json:"is_aktif"`
	DtmCrt          time.Time `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd          time.Time `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	SalesmanIDOld   int32     `gorm:"column:salesman_id_old" json:"salesman_id_old"`
	BranchIDOld     int32     `gorm:"column:branch_id_old" json:"branch_id_old"`
	AreaSrIDOld     int16     `gorm:"column:area_sr_id_old" json:"area_sr_id_old"`
	BranchIDNewOld  int16     `gorm:"column:branch_id_new_old" json:"branch_id_new_old"`
	DeviceID        string    `gorm:"column:device_id" json:"device_id"`
	IsMultipleLogin int16     `gorm:"column:is_multiple_login;not null" json:"is_multiple_login"`
	LevelID         string    `gorm:"column:level_id" json:"level_id"`
	SrID            string    `gorm:"column:sr_id" json:"sr_id"`
	RayonID         string    `gorm:"column:rayon_id" json:"rayon_id"`
	BranchID        string    `gorm:"column:branch_id" json:"branch_id"`
	AreaID          string    `gorm:"column:area_id" json:"area_id"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
