package xgorm

import (
	"errors"
	"gorm.io/gorm"
)

type Result struct {
	Error        error
	RowsAffected uint
}

func newResult(err error, rowsAffected uint) *Result {
	return &Result{
		Error:        err,
		RowsAffected: rowsAffected,
	}
}

func newResultByDB(db *gorm.DB) *Result {
	return newResult(db.Error, uint(db.RowsAffected))
}

func (r Result) Result() (uint, error) {
	return r.RowsAffected, r.Error
}

type CountResult struct {
	Error error
	Count uint
}

func newCountResult(err error, count uint) *CountResult {
	return &CountResult{
		Error: err,
		Count: count,
	}
}

func (r CountResult) Result() (uint, error) {
	return r.Count, r.Error
}

type IOperator interface {
	Create(md interface{}) *Result
	Creates(mds interface{}, batchSize int) *Result

	DeleteByID(id uint) *Result
	DeleteByQuery(condition *Condition) *Result

	UpdateByID(id uint, data map[string]interface{}) *Result
	UpdateByQuery(condition *Condition, data map[string]interface{}) *Result

	FindByID(id uint, destination interface{}, options ...FindByIDOption) error
	First(destination interface{}, options ...FindOneOption) error
	Last(destination interface{}, options ...FindOneOption) error
	Search(destination interface{}, options ...SearchOption) *CountResult
	Count(options ...CountOption) *CountResult
}

type Operator struct {
	*gorm.DB
	table interface{}
}

func NewOperator(db *gorm.DB, table interface{}) *Operator {
	return &Operator{
		DB:    db,
		table: table,
	}
}

func (o *Operator) Create(md interface{}) *Result {
	return newResultByDB(o.getDB().
		Create(md),
	)
}

func (o *Operator) Creates(mds interface{}, batchSize int) *Result {
	return newResultByDB(o.getDB().
		CreateInBatches(mds, batchSize),
	)
}

func (o *Operator) DeleteByID(id uint) *Result {
	return newResultByDB(o.getDB().
		Where("id = ?", id).Delete(nil),
	)
}

func (o *Operator) DeleteByQuery(condition *Condition) *Result {
	if condition == nil {
		return newResult(errors.New("delete by query, condition can't nil"), 0)
	}
	return newResultByDB(o.getDB().
		Where(condition.Query, condition.Args...).Delete(nil),
	)
}

func (o *Operator) UpdateByID(id uint, data map[string]interface{}) *Result {
	return newResultByDB(o.getDB().
		Where("id = ?", id).
		Updates(data),
	)
}

func (o *Operator) UpdateByQuery(condition *Condition, data map[string]interface{}) *Result {
	if condition == nil {
		return newResult(errors.New("update by query, condition can't nil"), 0)
	}
	return newResultByDB(o.getDB().
		Where(condition.Query, condition.Args...).
		Updates(data),
	)
}

func (o *Operator) FindByID(id uint, destination interface{}, options ...FindByIDOption) error {
	conf := newFindByIDOptions(options...)
	db := o.getDB()
	if conf.unscoped {
		db = db.Unscoped()
	}
	if len(conf.select_) > 0 {
		db = db.Select(conf.select_)
	}
	if len(conf.unselect) > 0 {
		db = db.Omit(conf.unselect...)
	}
	return db.
		Where("`id` = ?", id).
		First(destination).Error
}

func (o *Operator) First(destination interface{}, options ...FindOneOption) error {
	return o.getDBWithFindOneOptions(options...).
		First(destination).Error
}

func (o *Operator) Last(destination interface{}, options ...FindOneOption) error {
	return o.getDBWithFindOneOptions(options...).
		Last(destination).Error
}

func (o *Operator) Search(destination interface{}, options ...SearchOption) *CountResult {
	var count int64
	db := o.getDB()
	conf := newSearchOptions(options...)
	if conf.unscoped {
		db = db.Unscoped()
	}
	if len(conf.select_) > 0 {
		db = db.Select(conf.select_)
	}
	if len(conf.unselect) > 0 {
		db = db.Omit(conf.unselect...)
	}
	if conf.where != nil {
		db = db.Where(conf.where.Query, conf.where.Args...)
	}
	if conf.having != nil {
		db = db.Having(conf.having.Query, conf.having.Args...)
	}
	if conf.group != "" {
		db = db.Group(conf.group)
	}
	if conf.count {
		db = db.Count(&count)
	}
	if conf.limit != 0 {
		db = db.Limit(int(conf.limit))
	}
	if conf.page != 0 {
		db = db.Offset(int(conf.limit * (conf.page - 1)))
	}
	if conf.order != "" {
		db = db.Order(conf.order)
	}
	db = db.Find(destination)
	return newCountResult(db.Error, uint(count))
}

func (o *Operator) Count(options ...CountOption) *CountResult {
	var count int64
	db := o.getDB()
	conf := newCountOptions(options...)
	if conf.unscoped {
		db = db.Unscoped()
	}
	if conf.where != nil {
		db = db.Where(conf.where.Query, conf.where.Args...)
	}
	if conf.having != nil {
		db = db.Having(conf.having.Query, conf.having.Args...)
	}
	if conf.group != "" {
		db = db.Group(conf.group)
	}
	db = db.Count(&count)
	return newCountResult(db.Error, uint(count))
}

func (o *Operator) getDBWithFindOneOptions(options ...FindOneOption) *gorm.DB {
	db := o.getDB()
	conf := newFindOneOptions(options...)
	if conf.unscoped {
		db = db.Unscoped()
	}
	if len(conf.select_) > 0 {
		db = db.Select(conf.select_)
	}
	if len(conf.unselect) > 0 {
		db = db.Omit(conf.unselect...)
	}
	if conf.where != nil {
		db = db.Where(conf.where.Query, conf.where.Args...)
	}
	if conf.having != nil {
		db = db.Having(conf.having.Query, conf.having.Args...)
	}
	if conf.group != "" {
		db = db.Group(conf.group)
	}
	if conf.order != "" {
		db = db.Order(conf.order)
	}
	return db
}

func (o *Operator) getDB() *gorm.DB {
	switch o.table.(type) {
	case string:
		return o.DB.Table(o.table.(string))
	default:
		return o.DB.Model(o.table)
	}
}
