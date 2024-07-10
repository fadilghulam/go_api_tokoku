package model

import (
	"time"
)

const tableHrPerson = "hr_person"

type HrPerson struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	FullName         string    `gorm:"column:full_name;not null" json:"column:fullName"`
	Gender           string    `gorm:"column:gender;default:null" json:"column:gender"`
	Religion         string    `gorm:"column:religion;default:null" json:"column:religion"`
	BirthPlace       string    `gorm:"column:birth_place;default:null" json:"column:birthPlace"`
	BirthDate        time.Time `gorm:"column:birth_date;default:null" json:"column:birthDate"`
	MaritalStatus    string    `gorm:"column:marital_status;default:null" json:"column:maritalStatus"`
	Nationality      string    `gorm:"column:nationality;default:null" json:"column:nationality"`
	Address          string    `gorm:"column:address;default:null" json:"column:address"`
	Phone            string    `gorm:"column:phone;default:null" json:"column:phone"`
	Email            string    `gorm:"column:email;default:null" json:"column:email"`
	Ktp              string    `gorm:"column:ktp;default:null" json:"column:ktp"`
	KartuKeluarga    string    `gorm:"column:kartu_keluarga;default:null" json:"column:kartuKeluarga"`
	Photo            string    `gorm:"column:photo;default:null" json:"column:photo"`
	IsActive         int16     `gorm:"column:is_active;not null;default:1" json:"column:isActive"`
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:now()" json:"column:createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null;default:now()" json:"column:updatedAt"`
	Summary          string    `gorm:"column:summary;default:null" json:"column:summary"`
	FileKtp          string    `gorm:"column:file_ktp;default:null" json:"column:fileKtp"`
	Password         string    `gorm:"column:password;default:null" json:"column:password"`
	Domisili         string    `gorm:"column:domisili;default:null" json:"column:domisili"`
	Blood            string    `gorm:"column:blood;default:null" json:"column:blood"`
	Hobby            string    `gorm:"column:hobby;default:null" json:"column:hobby"`
	FileCv           string    `gorm:"column:file_cv;default:null" json:"column:fileCv"`
	KtpProvinsi      int32     `gorm:"column:ktp_provinsi;default:null" json:"column:ktpProvinsi"`
	KtpKabupaten     int32     `gorm:"column:ktp_kabupaten;default:null" json:"column:ktpKabupaten"`
	KtpKecamatan     int32     `gorm:"column:ktp_kecamatan;default:null" json:"column:ktpKecamatan"`
	KtpKelurahan     int64     `gorm:"column:ktp_kelurahan;default:null" json:"column:ktpKelurahan"`
	DmsProvinsi      int32     `gorm:"column:dms_provinsi;default:null" json:"column:dmsProvinsi"`
	DmsKabupaten     int32     `gorm:"column:dms_kabupaten;default:null" json:"column:dmsKabupaten"`
	DmsKecamatan     int32     `gorm:"column:dms_kecamatan;default:null" json:"column:dmsKecamatan"`
	DmsKelurahan     int64     `gorm:"column:dms_kelurahan;default:null" json:"column:dmsKelurahan"`
	KtpProvinsiName  string    `gorm:"column:ktp_provinsi_name;default:null" json:"column:ktpProvinsiName"`
	KtpKabupatenName string    `gorm:"column:ktp_kabupaten_name;default:null" json:"column:ktpKabupatenName"`
	KtpKecamatanName string    `gorm:"column:ktp_kecamatan_name;default:null" json:"column:ktpKecamatanName"`
	KtpKelurahanName string    `gorm:"column:ktp_kelurahan_name;default:null" json:"column:ktpKelurahanName"`
	DmsProvinsiName  string    `gorm:"column:dms_provinsi_name;default:null" json:"column:dmsProvinsiName"`
	DmsKabupatenName string    `gorm:"column:dms_kabupaten_name;default:null" json:"column:dmsKabupatenName"`
	DmsKecamatanName string    `gorm:"column:dms_kecamatan_name;default:null" json:"column:dmsKecamatanName"`
	DmsKelurahanName string    `gorm:"column:dms_kelurahan_name;default:null" json:"column:dmsKelurahanName"`
	DescHobby        string    `gorm:"column:desc_hobby;default:null" json:"column:descHobby"`
	Kelebihan        string    `gorm:"column:kelebihan;default:null" json:"column:kelebihan"`
	DescKelebihan    string    `gorm:"column:desc_kelebihan;default:null" json:"column:descKelebihan"`
	UserID           int32     `gorm:"column:user_id;default:null" json:"column:userId"`
	Kekurangan       string    `gorm:"column:kekurangan;default:null" json:"column:kekurangan"`
	DescKekurangan   string    `gorm:"column:desc_kekurangan;default:null" json:"column:descKekurangan"`
	JmlAnak          int16     `gorm:"column:jml_anak;default:null" json:"column:jmlAnak"`
	AppId            int32     `gorm:"column:app_id;default:17" json:"column:AppId"`
}

func (*HrPerson) TableName() string {
	return tableHrPerson
}
