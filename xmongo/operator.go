package xmongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type ObjectID = primitive.ObjectID
type D = primitive.D
type E = primitive.E

type CreateResult struct {
	Error      error
	InsertedID interface{}
}

type CreatesResult struct {
	Error       error
	InsertedIDs []interface{}
}

type UpdateResult struct {
	Error         error
	MatchedCount  int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID    interface{}
}

type FindOneResult struct {
	Error error
}

type CountResult struct {
	Error error
	Count uint
}

type IOperator interface {
	Create(md interface{}) CreateResult
	Creates(document []interface{}) CreatesResult

	DeleteByID(id uint) UpdateResult
	DeleteByQuery(condition map[string]interface{}) UpdateResult

	UpdateByID(id uint, data map[string]interface{}) UpdateResult
	UpdateByQuery(condition map[string]interface{}, data map[string]interface{}) UpdateResult

	FindByID(id ObjectID, destination interface{}, options ...FindByIDOption) error
	First(destination interface{}, options ...FindOneOption) error
	Search(destination interface{}, options ...SearchOption) CountResult
	Count(options ...CountOption) CountResult
}

type Operator struct {
	*mongo.Client
	Database   string
	Collection string
}

func (o *Operator) Create(md interface{}) CreateResult {
	err := o.mustStructPointer(md)
	if err != nil {
		return CreateResult{Error: err}
	}
	now := time.Now()
	o.setField(md, "CreatedAt", now)
	o.setField(md, "UpdatedAt", now)
	o.setField(md, "ID", primitive.NewObjectID())
	result, err := o.collection().InsertOne(context.TODO(), md)
	if err != nil {
		return CreateResult{Error: err}
	}
	return CreateResult{InsertedID: result.InsertedID}
}

func (o *Operator) Creates(mds []interface{}) CreatesResult {
	now := time.Now()
	for _, md := range mds {
		err := o.mustStructPointer(md)
		if err != nil {
			return CreatesResult{Error: err}
		}
		o.setField(md, "CreatedAt", now)
		o.setField(md, "UpdatedAt", now)
		o.setField(md, "ID", primitive.NewObjectID())
	}
	result, err := o.collection().InsertMany(context.TODO(), mds)
	if err != nil {
		return CreatesResult{Error: err}
	}
	return CreatesResult{InsertedIDs: result.InsertedIDs}
}

func (o *Operator) DeleteByID(id uint) UpdateResult {
	result, err := o.collection().UpdateOne(context.TODO(), map[string]interface{}{"_id": id}, o.deletedAt())
	if err != nil {
		return UpdateResult{Error: err}
	}
	return UpdateResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}
}

func (o *Operator) DeleteByQuery(query map[string]interface{}) UpdateResult {
	result, err := o.collection().UpdateMany(context.TODO(), query, o.deletedAt())
	if err != nil {
		return UpdateResult{Error: err}
	}
	return UpdateResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}
}

func (o *Operator) UpdateByID(id uint, data map[string]interface{}) UpdateResult {
	result, err := o.collection().UpdateOne(context.TODO(), map[string]interface{}{"_id": id}, o.updatedAt(data))
	if err != nil {
		return UpdateResult{Error: err}
	}
	return UpdateResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}
}

func (o *Operator) UpdateByQuery(condition map[string]interface{}, data map[string]interface{}) UpdateResult {
	result, err := o.collection().UpdateMany(context.TODO(), condition, o.updatedAt(data))
	if err != nil {
		return UpdateResult{Error: err}
	}
	return UpdateResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}
}

func (o *Operator) FindByID(id ObjectID, destination interface{}, options ...FindByIDOption) error {
	condition := map[string]interface{}{"_id": id}
	conf := newFindByIDOptions(options...)
	if !conf.unscoped {
		o.scoped(condition)
	}
	// mongo options
	opts := &mongoOptions.FindOneOptions{}
	projection := map[string]interface{}{}
	for _, item := range conf.select_ {
		projection[item] = 1
	}
	for _, item := range conf.unselect {
		projection[item] = -1
	}
	if len(projection) > 0 {
		opts.SetProjection(projection)
	}
	// exec
	return o.collection().FindOne(context.TODO(), condition, opts).Decode(destination)
}

func (o *Operator) First(destination interface{}, options ...FindOneOption) error {
	conf := newFindOneOptions(options...)
	if !conf.unscoped {
		o.scoped(conf.condition)
	}
	// mongo options
	opts := &mongoOptions.FindOneOptions{}
	if conf.order != nil {
		opts.SetSort(conf.order)
	}
	projection := map[string]interface{}{}
	for _, item := range conf.select_ {
		projection[item] = 1
	}
	for _, item := range conf.unselect {
		projection[item] = -1
	}
	if len(projection) > 0 {
		opts.SetProjection(projection)
	}
	// exec
	return o.collection().FindOne(context.TODO(), conf.condition, opts).Decode(destination)
}

func (o *Operator) Search(destination interface{}, options ...SearchOption) CountResult {
	conf := newSearchOptions(options...)
	if !conf.unscoped {
		o.scoped(conf.condition)
	}
	// mongo options
	opts := &mongoOptions.FindOptions{}
	if conf.limit != 0 {
		opts.SetLimit(int64(conf.limit))
	}
	if conf.page != 0 {
		opts.SetSkip(int64(conf.limit * (conf.page - 1)))
	}
	if conf.order != nil {
		opts.SetSort(conf.order)
	}
	projection := map[string]interface{}{}
	for _, item := range conf.select_ {
		projection[item] = 1
	}
	for _, item := range conf.unselect {
		projection[item] = -1
	}
	if len(projection) > 0 {
		opts.SetProjection(projection)
	}

	// exec search
	ctx := context.TODO()
	cursor, err := o.collection().Find(ctx, conf.condition, opts)
	if err != nil {
		return CountResult{Error: err}
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)
	err = cursor.All(ctx, destination)
	if err != nil {
		return CountResult{Error: err}
	}
	if !conf.count {
		return CountResult{}
	}
	count, err := o.collection().CountDocuments(context.TODO(), conf.condition)
	if err != nil {
		return CountResult{Error: err}
	}
	return CountResult{Count: uint(count)}
}

func (o *Operator) Count(options ...CountOption) CountResult {
	conf := newCountOptions(options...)
	if !conf.unscoped {
		o.scoped(conf.condition)
	}
	count, err := o.collection().CountDocuments(context.TODO(), conf.condition)
	if err != nil {
		return CountResult{Error: err}
	}
	return CountResult{Count: uint(count)}
}

func (o *Operator) deletedAt() map[string]interface{} {
	return map[string]interface{}{
		"$set": map[string]interface{}{
			"deleted_at": time.Now(),
		},
	}
}

func (o *Operator) updatedAt(update map[string]interface{}) map[string]interface{} {
	if update == nil {
		update = map[string]interface{}{}
	}
	if _, ok := update["$set"]; !ok {
		update["$set"] = map[string]interface{}{}
	}
	if v, ok := update["$set"].(map[string]interface{}); ok {
		if _, ok := v["updated_at"]; !ok {
			v["updated_at"] = time.Now()
		}
	}
	return update
}

func (o *Operator) scoped(filter map[string]interface{}) {
	if filter == nil {
		filter = map[string]interface{}{}
	}
	filter["deleted_at"] = map[string]interface{}{"$exists": false}
}

func (o *Operator) collection() *mongo.Collection {
	return o.Client.Database(o.Database).Collection(o.Collection)
}

func (o *Operator) mustStructPointer(v interface{}) error {
	if reflect.TypeOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).Elem().Kind() == reflect.Struct {
		return nil
	}
	return errors.New("document must be a struct pointer")
}

func (o *Operator) setField(md interface{}, key string, val interface{}) {
	elem := reflect.ValueOf(md).Elem()
	if field, ok := elem.Type().FieldByName(key); ok {
		if elem.FieldByIndex(field.Index).IsZero() {
			elem.FieldByName(key).Set(reflect.ValueOf(val))
		}
	}
}
