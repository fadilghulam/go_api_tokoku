// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package generated

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"go_api_tokoku/internal/model"
)

func newCustomer(db *gorm.DB, opts ...gen.DOOption) customer {
	_customer := customer{}

	_customer.customerDo.UseDB(db, opts...)
	_customer.customerDo.UseModel(&model.Customer{})

	tableName := _customer.customerDo.TableName()
	_customer.ALL = field.NewAsterisk(tableName)
	_customer.ID = field.NewInt64(tableName, "id")
	_customer.SalesmanID = field.NewInt32(tableName, "salesman_id")
	_customer.Name = field.NewString(tableName, "name")
	_customer.OutletName = field.NewString(tableName, "outlet_name")
	_customer.Alamat = field.NewString(tableName, "alamat")
	_customer.Phone = field.NewString(tableName, "phone")
	_customer.Tipe = field.NewInt32(tableName, "tipe")
	_customer.LatitudeLongitude = field.NewString(tableName, "latitude_longitude")
	_customer.ImageKtp = field.NewString(tableName, "image_ktp")
	_customer.ImageToko = field.NewString(tableName, "image_toko")
	_customer.IsAcc = field.NewInt32(tableName, "is_acc")
	_customer.IsAktif = field.NewInt32(tableName, "is_aktif")
	_customer.DtmCrt = field.NewTime(tableName, "dtm_crt")
	_customer.DtmUpd = field.NewTime(tableName, "dtm_upd")
	_customer.Plafond = field.NewFloat64(tableName, "plafond")
	_customer.Piutang = field.NewFloat64(tableName, "piutang")
	_customer.KodeCustomer = field.NewString(tableName, "kode_customer")
	_customer.QrCode = field.NewString(tableName, "qr_code")
	_customer.Nik = field.NewString(tableName, "nik")
	_customer.Diskon = field.NewFloat64(tableName, "diskon")
	_customer.Provinsi = field.NewString(tableName, "provinsi")
	_customer.Kabupaten = field.NewString(tableName, "kabupaten")
	_customer.Kecamatan = field.NewString(tableName, "kecamatan")
	_customer.Kelurahan = field.NewString(tableName, "kelurahan")
	_customer.KawasanToko = field.NewString(tableName, "kawasan_toko")
	_customer.HariKunjungan = field.NewString(tableName, "hari_kunjungan")
	_customer.FrekKunjungan = field.NewString(tableName, "frek_kunjungan")
	_customer.KawasanTokoOth = field.NewString(tableName, "kawasan_toko_oth")
	_customer.FrekKunjunganOth = field.NewString(tableName, "frek_kunjungan_oth")
	_customer.IsVerifikasi = field.NewInt32(tableName, "is_verifikasi")
	_customer.ImageTokoAfter = field.NewString(tableName, "image_toko_after")
	_customer.Validated = field.NewInt32(tableName, "validated")
	_customer.SyncKey = field.NewString(tableName, "sync_key")
	_customer.AksesDoubleKredit = field.NewTime(tableName, "akses_double_kredit")
	_customer.TanggalVerifikasi = field.NewTime(tableName, "tanggal_verifikasi")
	_customer.VerifiedBy = field.NewInt32(tableName, "verified_by")
	_customer.IsVerifikasiLokasi = field.NewInt16(tableName, "is_verifikasi_lokasi")
	_customer.VisitExtra = field.NewTime(tableName, "visit_extra")
	_customer.IsMandiri = field.NewInt16(tableName, "is_mandiri")
	_customer.SetMandiriBy = field.NewInt32(tableName, "set_mandiri_by")
	_customer.IsKasus = field.NewInt16(tableName, "is_kasus")
	_customer.SalesmanTemp = field.NewInt64(tableName, "salesman_temp")
	_customer.SisaKreditNoo = field.NewInt16(tableName, "sisa_kredit_noo")
	_customer.AreaID = field.NewInt32(tableName, "area_id")
	_customer.SalesmanTypeID = field.NewInt16(tableName, "salesman_type_id")
	_customer.IsHandover = field.NewInt16(tableName, "is_handover")
	_customer.CreatedID = field.NewInt32(tableName, "created_id")
	_customer.DateLastTransaction = field.NewTime(tableName, "date_last_transaction")
	_customer.DateLastVisitBySalesman = field.NewTime(tableName, "date_last_visit_by_salesman")
	_customer.OutletPhoto = field.NewString(tableName, "outlet_photo")
	_customer.Note = field.NewString(tableName, "note")
	_customer.SubjectTypeID = field.NewInt16(tableName, "subject_type_id")
	_customer.SalesmanIDCreator = field.NewInt32(tableName, "salesman_id_creator")
	_customer.MerchandiserIDCreator = field.NewInt32(tableName, "merchandiser_id_creator")
	_customer.TeamleaderIDCreator = field.NewInt32(tableName, "teamleader_id_creator")
	_customer.BranchID = field.NewInt16(tableName, "branch_id")
	_customer.RayonID = field.NewInt16(tableName, "rayon_id")
	_customer.SrID = field.NewInt16(tableName, "sr_id")
	_customer.IsConsume = field.NewInt16(tableName, "is_consume")
	_customer.EmployeeIDConsume = field.NewInt64(tableName, "employee_id_consume")
	_customer.ConsumeAt = field.NewTime(tableName, "consume_at")
	_customer.Tag = field.NewString(tableName, "tag")

	_customer.fillFieldMap()

	return _customer
}

type customer struct {
	customerDo

	ALL                     field.Asterisk
	ID                      field.Int64
	SalesmanID              field.Int32
	Name                    field.String
	OutletName              field.String
	Alamat                  field.String
	Phone                   field.String
	Tipe                    field.Int32
	LatitudeLongitude       field.String
	ImageKtp                field.String
	ImageToko               field.String
	IsAcc                   field.Int32
	IsAktif                 field.Int32
	DtmCrt                  field.Time
	DtmUpd                  field.Time
	Plafond                 field.Float64
	Piutang                 field.Float64
	KodeCustomer            field.String
	QrCode                  field.String
	Nik                     field.String
	Diskon                  field.Float64
	Provinsi                field.String
	Kabupaten               field.String
	Kecamatan               field.String
	Kelurahan               field.String
	KawasanToko             field.String
	HariKunjungan           field.String
	FrekKunjungan           field.String
	KawasanTokoOth          field.String
	FrekKunjunganOth        field.String
	IsVerifikasi            field.Int32 // digunakan untuk membedakan new outlet disetujui
	ImageTokoAfter          field.String
	Validated               field.Int32
	SyncKey                 field.String
	AksesDoubleKredit       field.Time
	TanggalVerifikasi       field.Time
	VerifiedBy              field.Int32
	IsVerifikasiLokasi      field.Int16
	VisitExtra              field.Time
	IsMandiri               field.Int16
	SetMandiriBy            field.Int32 // diset oleh user id pada tabel user
	IsKasus                 field.Int16 // customer kasus...
	SalesmanTemp            field.Int64
	SisaKreditNoo           field.Int16
	AreaID                  field.Int32
	SalesmanTypeID          field.Int16
	IsHandover              field.Int16
	CreatedID               field.Int32 // employee_id
	DateLastTransaction     field.Time
	DateLastVisitBySalesman field.Time
	OutletPhoto             field.String
	Note                    field.String
	SubjectTypeID           field.Int16
	SalesmanIDCreator       field.Int32
	MerchandiserIDCreator   field.Int32
	TeamleaderIDCreator     field.Int32
	BranchID                field.Int16
	RayonID                 field.Int16
	SrID                    field.Int16
	IsConsume               field.Int16
	EmployeeIDConsume       field.Int64
	ConsumeAt               field.Time
	/*
		REGULAR
		RECOMENDATION
	*/
	Tag field.String

	fieldMap map[string]field.Expr
}

func (c customer) Table(newTableName string) *customer {
	c.customerDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c customer) As(alias string) *customer {
	c.customerDo.DO = *(c.customerDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *customer) updateTableName(table string) *customer {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewInt64(table, "id")
	c.SalesmanID = field.NewInt32(table, "salesman_id")
	c.Name = field.NewString(table, "name")
	c.OutletName = field.NewString(table, "outlet_name")
	c.Alamat = field.NewString(table, "alamat")
	c.Phone = field.NewString(table, "phone")
	c.Tipe = field.NewInt32(table, "tipe")
	c.LatitudeLongitude = field.NewString(table, "latitude_longitude")
	c.ImageKtp = field.NewString(table, "image_ktp")
	c.ImageToko = field.NewString(table, "image_toko")
	c.IsAcc = field.NewInt32(table, "is_acc")
	c.IsAktif = field.NewInt32(table, "is_aktif")
	c.DtmCrt = field.NewTime(table, "dtm_crt")
	c.DtmUpd = field.NewTime(table, "dtm_upd")
	c.Plafond = field.NewFloat64(table, "plafond")
	c.Piutang = field.NewFloat64(table, "piutang")
	c.KodeCustomer = field.NewString(table, "kode_customer")
	c.QrCode = field.NewString(table, "qr_code")
	c.Nik = field.NewString(table, "nik")
	c.Diskon = field.NewFloat64(table, "diskon")
	c.Provinsi = field.NewString(table, "provinsi")
	c.Kabupaten = field.NewString(table, "kabupaten")
	c.Kecamatan = field.NewString(table, "kecamatan")
	c.Kelurahan = field.NewString(table, "kelurahan")
	c.KawasanToko = field.NewString(table, "kawasan_toko")
	c.HariKunjungan = field.NewString(table, "hari_kunjungan")
	c.FrekKunjungan = field.NewString(table, "frek_kunjungan")
	c.KawasanTokoOth = field.NewString(table, "kawasan_toko_oth")
	c.FrekKunjunganOth = field.NewString(table, "frek_kunjungan_oth")
	c.IsVerifikasi = field.NewInt32(table, "is_verifikasi")
	c.ImageTokoAfter = field.NewString(table, "image_toko_after")
	c.Validated = field.NewInt32(table, "validated")
	c.SyncKey = field.NewString(table, "sync_key")
	c.AksesDoubleKredit = field.NewTime(table, "akses_double_kredit")
	c.TanggalVerifikasi = field.NewTime(table, "tanggal_verifikasi")
	c.VerifiedBy = field.NewInt32(table, "verified_by")
	c.IsVerifikasiLokasi = field.NewInt16(table, "is_verifikasi_lokasi")
	c.VisitExtra = field.NewTime(table, "visit_extra")
	c.IsMandiri = field.NewInt16(table, "is_mandiri")
	c.SetMandiriBy = field.NewInt32(table, "set_mandiri_by")
	c.IsKasus = field.NewInt16(table, "is_kasus")
	c.SalesmanTemp = field.NewInt64(table, "salesman_temp")
	c.SisaKreditNoo = field.NewInt16(table, "sisa_kredit_noo")
	c.AreaID = field.NewInt32(table, "area_id")
	c.SalesmanTypeID = field.NewInt16(table, "salesman_type_id")
	c.IsHandover = field.NewInt16(table, "is_handover")
	c.CreatedID = field.NewInt32(table, "created_id")
	c.DateLastTransaction = field.NewTime(table, "date_last_transaction")
	c.DateLastVisitBySalesman = field.NewTime(table, "date_last_visit_by_salesman")
	c.OutletPhoto = field.NewString(table, "outlet_photo")
	c.Note = field.NewString(table, "note")
	c.SubjectTypeID = field.NewInt16(table, "subject_type_id")
	c.SalesmanIDCreator = field.NewInt32(table, "salesman_id_creator")
	c.MerchandiserIDCreator = field.NewInt32(table, "merchandiser_id_creator")
	c.TeamleaderIDCreator = field.NewInt32(table, "teamleader_id_creator")
	c.BranchID = field.NewInt16(table, "branch_id")
	c.RayonID = field.NewInt16(table, "rayon_id")
	c.SrID = field.NewInt16(table, "sr_id")
	c.IsConsume = field.NewInt16(table, "is_consume")
	c.EmployeeIDConsume = field.NewInt64(table, "employee_id_consume")
	c.ConsumeAt = field.NewTime(table, "consume_at")
	c.Tag = field.NewString(table, "tag")

	c.fillFieldMap()

	return c
}

func (c *customer) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *customer) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 62)
	c.fieldMap["id"] = c.ID
	c.fieldMap["salesman_id"] = c.SalesmanID
	c.fieldMap["name"] = c.Name
	c.fieldMap["outlet_name"] = c.OutletName
	c.fieldMap["alamat"] = c.Alamat
	c.fieldMap["phone"] = c.Phone
	c.fieldMap["tipe"] = c.Tipe
	c.fieldMap["latitude_longitude"] = c.LatitudeLongitude
	c.fieldMap["image_ktp"] = c.ImageKtp
	c.fieldMap["image_toko"] = c.ImageToko
	c.fieldMap["is_acc"] = c.IsAcc
	c.fieldMap["is_aktif"] = c.IsAktif
	c.fieldMap["dtm_crt"] = c.DtmCrt
	c.fieldMap["dtm_upd"] = c.DtmUpd
	c.fieldMap["plafond"] = c.Plafond
	c.fieldMap["piutang"] = c.Piutang
	c.fieldMap["kode_customer"] = c.KodeCustomer
	c.fieldMap["qr_code"] = c.QrCode
	c.fieldMap["nik"] = c.Nik
	c.fieldMap["diskon"] = c.Diskon
	c.fieldMap["provinsi"] = c.Provinsi
	c.fieldMap["kabupaten"] = c.Kabupaten
	c.fieldMap["kecamatan"] = c.Kecamatan
	c.fieldMap["kelurahan"] = c.Kelurahan
	c.fieldMap["kawasan_toko"] = c.KawasanToko
	c.fieldMap["hari_kunjungan"] = c.HariKunjungan
	c.fieldMap["frek_kunjungan"] = c.FrekKunjungan
	c.fieldMap["kawasan_toko_oth"] = c.KawasanTokoOth
	c.fieldMap["frek_kunjungan_oth"] = c.FrekKunjunganOth
	c.fieldMap["is_verifikasi"] = c.IsVerifikasi
	c.fieldMap["image_toko_after"] = c.ImageTokoAfter
	c.fieldMap["validated"] = c.Validated
	c.fieldMap["sync_key"] = c.SyncKey
	c.fieldMap["akses_double_kredit"] = c.AksesDoubleKredit
	c.fieldMap["tanggal_verifikasi"] = c.TanggalVerifikasi
	c.fieldMap["verified_by"] = c.VerifiedBy
	c.fieldMap["is_verifikasi_lokasi"] = c.IsVerifikasiLokasi
	c.fieldMap["visit_extra"] = c.VisitExtra
	c.fieldMap["is_mandiri"] = c.IsMandiri
	c.fieldMap["set_mandiri_by"] = c.SetMandiriBy
	c.fieldMap["is_kasus"] = c.IsKasus
	c.fieldMap["salesman_temp"] = c.SalesmanTemp
	c.fieldMap["sisa_kredit_noo"] = c.SisaKreditNoo
	c.fieldMap["area_id"] = c.AreaID
	c.fieldMap["salesman_type_id"] = c.SalesmanTypeID
	c.fieldMap["is_handover"] = c.IsHandover
	c.fieldMap["created_id"] = c.CreatedID
	c.fieldMap["date_last_transaction"] = c.DateLastTransaction
	c.fieldMap["date_last_visit_by_salesman"] = c.DateLastVisitBySalesman
	c.fieldMap["outlet_photo"] = c.OutletPhoto
	c.fieldMap["note"] = c.Note
	c.fieldMap["subject_type_id"] = c.SubjectTypeID
	c.fieldMap["salesman_id_creator"] = c.SalesmanIDCreator
	c.fieldMap["merchandiser_id_creator"] = c.MerchandiserIDCreator
	c.fieldMap["teamleader_id_creator"] = c.TeamleaderIDCreator
	c.fieldMap["branch_id"] = c.BranchID
	c.fieldMap["rayon_id"] = c.RayonID
	c.fieldMap["sr_id"] = c.SrID
	c.fieldMap["is_consume"] = c.IsConsume
	c.fieldMap["employee_id_consume"] = c.EmployeeIDConsume
	c.fieldMap["consume_at"] = c.ConsumeAt
	c.fieldMap["tag"] = c.Tag
}

func (c customer) clone(db *gorm.DB) customer {
	c.customerDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c customer) replaceDB(db *gorm.DB) customer {
	c.customerDo.ReplaceDB(db)
	return c
}

type customerDo struct{ gen.DO }

type ICustomerDo interface {
	gen.SubQuery
	Debug() ICustomerDo
	WithContext(ctx context.Context) ICustomerDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ICustomerDo
	WriteDB() ICustomerDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ICustomerDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ICustomerDo
	Not(conds ...gen.Condition) ICustomerDo
	Or(conds ...gen.Condition) ICustomerDo
	Select(conds ...field.Expr) ICustomerDo
	Where(conds ...gen.Condition) ICustomerDo
	Order(conds ...field.Expr) ICustomerDo
	Distinct(cols ...field.Expr) ICustomerDo
	Omit(cols ...field.Expr) ICustomerDo
	Join(table schema.Tabler, on ...field.Expr) ICustomerDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ICustomerDo
	RightJoin(table schema.Tabler, on ...field.Expr) ICustomerDo
	Group(cols ...field.Expr) ICustomerDo
	Having(conds ...gen.Condition) ICustomerDo
	Limit(limit int) ICustomerDo
	Offset(offset int) ICustomerDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ICustomerDo
	Unscoped() ICustomerDo
	Create(values ...*model.Customer) error
	CreateInBatches(values []*model.Customer, batchSize int) error
	Save(values ...*model.Customer) error
	First() (*model.Customer, error)
	Take() (*model.Customer, error)
	Last() (*model.Customer, error)
	Find() ([]*model.Customer, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Customer, err error)
	FindInBatches(result *[]*model.Customer, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Customer) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ICustomerDo
	Assign(attrs ...field.AssignExpr) ICustomerDo
	Joins(fields ...field.RelationField) ICustomerDo
	Preload(fields ...field.RelationField) ICustomerDo
	FirstOrInit() (*model.Customer, error)
	FirstOrCreate() (*model.Customer, error)
	FindByPage(offset int, limit int) (result []*model.Customer, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ICustomerDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (c customerDo) Debug() ICustomerDo {
	return c.withDO(c.DO.Debug())
}

func (c customerDo) WithContext(ctx context.Context) ICustomerDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c customerDo) ReadDB() ICustomerDo {
	return c.Clauses(dbresolver.Read)
}

func (c customerDo) WriteDB() ICustomerDo {
	return c.Clauses(dbresolver.Write)
}

func (c customerDo) Session(config *gorm.Session) ICustomerDo {
	return c.withDO(c.DO.Session(config))
}

func (c customerDo) Clauses(conds ...clause.Expression) ICustomerDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c customerDo) Returning(value interface{}, columns ...string) ICustomerDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c customerDo) Not(conds ...gen.Condition) ICustomerDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c customerDo) Or(conds ...gen.Condition) ICustomerDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c customerDo) Select(conds ...field.Expr) ICustomerDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c customerDo) Where(conds ...gen.Condition) ICustomerDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c customerDo) Order(conds ...field.Expr) ICustomerDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c customerDo) Distinct(cols ...field.Expr) ICustomerDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c customerDo) Omit(cols ...field.Expr) ICustomerDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c customerDo) Join(table schema.Tabler, on ...field.Expr) ICustomerDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c customerDo) LeftJoin(table schema.Tabler, on ...field.Expr) ICustomerDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c customerDo) RightJoin(table schema.Tabler, on ...field.Expr) ICustomerDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c customerDo) Group(cols ...field.Expr) ICustomerDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c customerDo) Having(conds ...gen.Condition) ICustomerDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c customerDo) Limit(limit int) ICustomerDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c customerDo) Offset(offset int) ICustomerDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c customerDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ICustomerDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c customerDo) Unscoped() ICustomerDo {
	return c.withDO(c.DO.Unscoped())
}

func (c customerDo) Create(values ...*model.Customer) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c customerDo) CreateInBatches(values []*model.Customer, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c customerDo) Save(values ...*model.Customer) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c customerDo) First() (*model.Customer, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) Take() (*model.Customer, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) Last() (*model.Customer, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) Find() ([]*model.Customer, error) {
	result, err := c.DO.Find()
	return result.([]*model.Customer), err
}

func (c customerDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Customer, err error) {
	buf := make([]*model.Customer, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c customerDo) FindInBatches(result *[]*model.Customer, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c customerDo) Attrs(attrs ...field.AssignExpr) ICustomerDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c customerDo) Assign(attrs ...field.AssignExpr) ICustomerDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c customerDo) Joins(fields ...field.RelationField) ICustomerDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c customerDo) Preload(fields ...field.RelationField) ICustomerDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c customerDo) FirstOrInit() (*model.Customer, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) FirstOrCreate() (*model.Customer, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) FindByPage(offset int, limit int) (result []*model.Customer, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c customerDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c customerDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c customerDo) Delete(models ...*model.Customer) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *customerDo) withDO(do gen.Dao) *customerDo {
	c.DO = *do.(*gen.DO)
	return c
}
