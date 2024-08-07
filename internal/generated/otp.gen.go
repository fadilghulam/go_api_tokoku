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

func newOtp(db *gorm.DB, opts ...gen.DOOption) otp {
	_otp := otp{}

	_otp.otpDo.UseDB(db, opts...)
	_otp.otpDo.UseModel(&model.Otp{})

	tableName := _otp.otpDo.TableName()
	_otp.ALL = field.NewAsterisk(tableName)
	_otp.ID = field.NewInt32(tableName, "id")
	_otp.AppName = field.NewString(tableName, "app_name")
	_otp.Description = field.NewString(tableName, "description")
	_otp.Type = field.NewString(tableName, "type")
	_otp.SendTo = field.NewString(tableName, "send_to")
	_otp.Otp = field.NewString(tableName, "otp")
	_otp.UserID = field.NewInt32(tableName, "user_id")
	_otp.ExpiredAt = field.NewTime(tableName, "expired_at")
	_otp.ConfirmedAt = field.NewTime(tableName, "confirmed_at")
	_otp.CreatedAt = field.NewTime(tableName, "created_at")
	_otp.UpdatedAt = field.NewTime(tableName, "updated_at")
	_otp.DeletedAt = field.NewField(tableName, "deleted_at")
	_otp.Label = field.NewString(tableName, "label")

	_otp.fillFieldMap()

	return _otp
}

type otp struct {
	otpDo

	ALL         field.Asterisk
	ID          field.Int32
	AppName     field.String
	Description field.String
	Type        field.String
	SendTo      field.String
	Otp         field.String
	UserID      field.Int32
	ExpiredAt   field.Time
	ConfirmedAt field.Time
	CreatedAt   field.Time
	UpdatedAt   field.Time
	DeletedAt   field.Field
	Label       field.String

	fieldMap map[string]field.Expr
}

func (o otp) Table(newTableName string) *otp {
	o.otpDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o otp) As(alias string) *otp {
	o.otpDo.DO = *(o.otpDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *otp) updateTableName(table string) *otp {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewInt32(table, "id")
	o.AppName = field.NewString(table, "app_name")
	o.Description = field.NewString(table, "description")
	o.Type = field.NewString(table, "type")
	o.SendTo = field.NewString(table, "send_to")
	o.Otp = field.NewString(table, "otp")
	o.UserID = field.NewInt32(table, "user_id")
	o.ExpiredAt = field.NewTime(table, "expired_at")
	o.ConfirmedAt = field.NewTime(table, "confirmed_at")
	o.CreatedAt = field.NewTime(table, "created_at")
	o.UpdatedAt = field.NewTime(table, "updated_at")
	o.DeletedAt = field.NewField(table, "deleted_at")
	o.Label = field.NewString(table, "label")

	o.fillFieldMap()

	return o
}

func (o *otp) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *otp) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 13)
	o.fieldMap["id"] = o.ID
	o.fieldMap["app_name"] = o.AppName
	o.fieldMap["description"] = o.Description
	o.fieldMap["type"] = o.Type
	o.fieldMap["send_to"] = o.SendTo
	o.fieldMap["otp"] = o.Otp
	o.fieldMap["user_id"] = o.UserID
	o.fieldMap["expired_at"] = o.ExpiredAt
	o.fieldMap["confirmed_at"] = o.ConfirmedAt
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
	o.fieldMap["deleted_at"] = o.DeletedAt
	o.fieldMap["label"] = o.Label
}

func (o otp) clone(db *gorm.DB) otp {
	o.otpDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o otp) replaceDB(db *gorm.DB) otp {
	o.otpDo.ReplaceDB(db)
	return o
}

type otpDo struct{ gen.DO }

type IOtpDo interface {
	gen.SubQuery
	Debug() IOtpDo
	WithContext(ctx context.Context) IOtpDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOtpDo
	WriteDB() IOtpDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOtpDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOtpDo
	Not(conds ...gen.Condition) IOtpDo
	Or(conds ...gen.Condition) IOtpDo
	Select(conds ...field.Expr) IOtpDo
	Where(conds ...gen.Condition) IOtpDo
	Order(conds ...field.Expr) IOtpDo
	Distinct(cols ...field.Expr) IOtpDo
	Omit(cols ...field.Expr) IOtpDo
	Join(table schema.Tabler, on ...field.Expr) IOtpDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOtpDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOtpDo
	Group(cols ...field.Expr) IOtpDo
	Having(conds ...gen.Condition) IOtpDo
	Limit(limit int) IOtpDo
	Offset(offset int) IOtpDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOtpDo
	Unscoped() IOtpDo
	Create(values ...*model.Otp) error
	CreateInBatches(values []*model.Otp, batchSize int) error
	Save(values ...*model.Otp) error
	First() (*model.Otp, error)
	Take() (*model.Otp, error)
	Last() (*model.Otp, error)
	Find() ([]*model.Otp, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Otp, err error)
	FindInBatches(result *[]*model.Otp, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Otp) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOtpDo
	Assign(attrs ...field.AssignExpr) IOtpDo
	Joins(fields ...field.RelationField) IOtpDo
	Preload(fields ...field.RelationField) IOtpDo
	FirstOrInit() (*model.Otp, error)
	FirstOrCreate() (*model.Otp, error)
	FindByPage(offset int, limit int) (result []*model.Otp, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOtpDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o otpDo) Debug() IOtpDo {
	return o.withDO(o.DO.Debug())
}

func (o otpDo) WithContext(ctx context.Context) IOtpDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o otpDo) ReadDB() IOtpDo {
	return o.Clauses(dbresolver.Read)
}

func (o otpDo) WriteDB() IOtpDo {
	return o.Clauses(dbresolver.Write)
}

func (o otpDo) Session(config *gorm.Session) IOtpDo {
	return o.withDO(o.DO.Session(config))
}

func (o otpDo) Clauses(conds ...clause.Expression) IOtpDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o otpDo) Returning(value interface{}, columns ...string) IOtpDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o otpDo) Not(conds ...gen.Condition) IOtpDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o otpDo) Or(conds ...gen.Condition) IOtpDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o otpDo) Select(conds ...field.Expr) IOtpDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o otpDo) Where(conds ...gen.Condition) IOtpDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o otpDo) Order(conds ...field.Expr) IOtpDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o otpDo) Distinct(cols ...field.Expr) IOtpDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o otpDo) Omit(cols ...field.Expr) IOtpDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o otpDo) Join(table schema.Tabler, on ...field.Expr) IOtpDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o otpDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOtpDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o otpDo) RightJoin(table schema.Tabler, on ...field.Expr) IOtpDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o otpDo) Group(cols ...field.Expr) IOtpDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o otpDo) Having(conds ...gen.Condition) IOtpDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o otpDo) Limit(limit int) IOtpDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o otpDo) Offset(offset int) IOtpDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o otpDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOtpDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o otpDo) Unscoped() IOtpDo {
	return o.withDO(o.DO.Unscoped())
}

func (o otpDo) Create(values ...*model.Otp) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o otpDo) CreateInBatches(values []*model.Otp, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o otpDo) Save(values ...*model.Otp) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o otpDo) First() (*model.Otp, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Otp), nil
	}
}

func (o otpDo) Take() (*model.Otp, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Otp), nil
	}
}

func (o otpDo) Last() (*model.Otp, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Otp), nil
	}
}

func (o otpDo) Find() ([]*model.Otp, error) {
	result, err := o.DO.Find()
	return result.([]*model.Otp), err
}

func (o otpDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Otp, err error) {
	buf := make([]*model.Otp, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o otpDo) FindInBatches(result *[]*model.Otp, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o otpDo) Attrs(attrs ...field.AssignExpr) IOtpDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o otpDo) Assign(attrs ...field.AssignExpr) IOtpDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o otpDo) Joins(fields ...field.RelationField) IOtpDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o otpDo) Preload(fields ...field.RelationField) IOtpDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o otpDo) FirstOrInit() (*model.Otp, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Otp), nil
	}
}

func (o otpDo) FirstOrCreate() (*model.Otp, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Otp), nil
	}
}

func (o otpDo) FindByPage(offset int, limit int) (result []*model.Otp, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o otpDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o otpDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o otpDo) Delete(models ...*model.Otp) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *otpDo) withDO(do gen.Dao) *otpDo {
	o.DO = *do.(*gen.DO)
	return o
}
