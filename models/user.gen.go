package model

import (
	"time"
)

const TableNameUser = "public.user"

// User mapped from table <user>
type User struct {
	ID              int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	FullName        string     `gorm:"column:full_name;not null" json:"full_name"`
	Username        string     `gorm:"column:username;not null" json:"username"`
	Password        string     `gorm:"column:password;not null" json:"password"`
	Images          string     `gorm:"column:images;default:null" json:"images"`
	LevelIDOld      int32      `gorm:"column:level_id_old;default:null" json:"level_id_old"`
	IsAktif         string     `gorm:"column:is_aktif;not null" json:"is_aktif"`
	DtmCrt          time.Time  `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd          time.Time  `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	SalesmanIDOld   int32      `gorm:"column:salesman_id_old;default:null" json:"salesman_id_old"`
	BranchIDOld     int32      `gorm:"column:branch_id_old;default:null" json:"branch_id_old"`
	AreaSrIDOld     int16      `gorm:"column:area_sr_id_old;default:null" json:"area_sr_id_old"`
	BranchIDNewOld  int16      `gorm:"column:branch_id_new_old;default:null" json:"branch_id_new_old"`
	DeviceID        string     `gorm:"column:device_id;default:null" json:"device_id"`
	IsMultipleLogin int16      `gorm:"column:is_multiple_login;not null" json:"is_multiple_login"`
	LevelID         Int32Array `gorm:"column:level_id;type:int[];default:null" json:"level_id"`
	SrID            Int32Array `gorm:"type:integer[];column:sr_id;default:null" json:"sr_id"`
	RayonID         Int32Array `gorm:"type:integer[];column:rayon_id;default:null" json:"rayon_id"`
	BranchID        Int32Array `gorm:"type:integer[];column:branch_id;default:null" json:"branch_id"`
	AreaID          Int32Array `gorm:"type:integer[];column:area_id;default:null" json:"area_id"`
	ProfilePhoto    string     `gorm:"column:profile_photo;default:null" json:"profilePhoto"`
	GoogleId        string     `gorm:"column:google_id;default:null" json:"googleId"`
	IsVerified      int16      `gorm:"column:is_verified;default:null" json:"isVerified"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
