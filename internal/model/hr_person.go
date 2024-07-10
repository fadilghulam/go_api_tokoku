package model

import (
	"database/sql"
	"time"
)

type HrPerson struct {
	ID               int            `db:"id"`
	FullName         string         `db:"full_name"`
	Gender           sql.NullString `db:"gender"`
	Religion         sql.NullString `db:"religion"`
	BirthPlace       sql.NullString `db:"birth_place"`
	BirthDate        sql.NullTime   `db:"birth_date"`
	MaritalStatus    sql.NullString `db:"marital_status"`
	Nationality      sql.NullString `db:"nationality"`
	Address          sql.NullString `db:"address"`
	Phone            sql.NullString `db:"phone"`
	Email            sql.NullString `db:"email"`
	Ktp              sql.NullString `db:"ktp"`
	KartuKeluarga    sql.NullString `db:"kartu_keluarga"`
	Photo            sql.NullString `db:"photo"`
	IsActive         int            `db:"is_active"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
	Summary          sql.NullString `db:"summary"`
	FileKtp          sql.NullString `db:"file_ktp"`
	Password         sql.NullString `db:"password"`
	Domisili         sql.NullString `db:"domisili"`
	Blood            sql.NullString `db:"blood"`
	Hobby            sql.NullString `db:"hobby"`
	FileCv           sql.NullString `db:"file_cv"`
	KtpProvinsi      sql.NullInt64  `db:"ktp_provinsi"`
	KtpKabupaten     sql.NullInt64  `db:"ktp_kabupaten"`
	KtpKecamatan     sql.NullInt64  `db:"ktp_kecamatan"`
	KtpKelurahan     sql.NullInt64  `db:"ktp_kelurahan"`
	DmsProvinsi      sql.NullInt64  `db:"dms_provinsi"`
	DmsKabupaten     sql.NullInt64  `db:"dms_kabupaten"`
	DmsKecamatan     sql.NullInt64  `db:"dms_kecamatan"`
	DmsKelurahan     sql.NullInt64  `db:"dms_kelurahan"`
	KtpProvinsiName  sql.NullString `db:"ktp_provinsi_name"`
	KtpKabupatenName sql.NullString `db:"ktp_kabupaten_name"`
	KtpKecamatanName sql.NullString `db:"ktp_kecamatan_name"`
	KtpKelurahanName sql.NullString `db:"ktp_kelurahan_name"`
	DmsProvinsiName  sql.NullString `db:"dms_provinsi_name"`
	DmsKabupatenName sql.NullString `db:"dms_kabupaten_name"`
	DmsKecamatanName sql.NullString `db:"dms_kecamatan_name"`
	DmsKelurahanName sql.NullString `db:"dms_kelurahan_name"`
	DescHobby        sql.NullString `db:"desc_hobby"`
	Kelebihan        sql.NullString `db:"kelebihan"`
	DescKelebihan    sql.NullString `db:"desc_kelebihan"`
	UserID           sql.NullInt64  `db:"user_id"`
	Kekurangan       sql.NullString `db:"kekurangan"`
	DescKekurangan   sql.NullString `db:"desc_kekurangan"`
	JmlAnak          sql.NullInt64  `db:"jml_anak"`
}
