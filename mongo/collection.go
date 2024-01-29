package mongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	collection *mongo.Collection

	registry *bsoncodec.Registry
}

// DropCollection drops collection
// it's safe even collection is not exists
func (c *Collection) DropCollection(ctx context.Context) error {
	return c.collection.Drop(ctx)
}

// CloneCollection creates a copy of the Collection
func (c *Collection) CloneCollection() (*mongo.Collection, error) {
	return c.collection.Clone()
}

func (c *Collection) GetCollectionName() string {
	return c.collection.Name()
}

func (c *Collection) Aggregate(ctx context.Context, pipe any, opts ...*options.AggregateOptions) *Cursor {
	cur, err := c.collection.Aggregate(ctx, pipe, opts...)
	return &Cursor{
		ctx:    ctx,
		cursor: cur,
		err:    err,
	}
}

func (c *Collection) Find(ctx context.Context, query any, opts ...*options.FindOptions) *Cursor {
	cur, err := c.collection.Find(ctx, query, opts...)
	return &Cursor{
		ctx:    ctx,
		cursor: cur,
		err:    err,
	}
}

func (c *Collection) Count(ctx context.Context) (int, error) {
	num, err := c.collection.CountDocuments(ctx, struct{}{})
	return int(num), err
}

func (c *Collection) FindOne(ctx context.Context, query any, opts ...*options.FindOneOptions) *SingleResult {
	return c.collection.FindOne(ctx, query, opts...)
}

func (c *Collection) FindOneAndDelete(ctx context.Context, query any, opts ...*options.FindOneAndDeleteOptions) *SingleResult {
	return c.collection.FindOneAndDelete(ctx, query, opts...)
}

func (c *Collection) FindOneAndReplace(ctx context.Context, query, doc any, opts ...*options.FindOneAndReplaceOptions) *SingleResult {
	return c.collection.FindOneAndReplace(ctx, query, doc, opts...)
}

func (c *Collection) FindOneAndUpdate(ctx context.Context, query, update any, opts ...*options.FindOneAndUpdateOptions) *SingleResult {
	return c.collection.FindOneAndUpdate(ctx, query, update, opts...)
}

func (c *Collection) InsertMany(ctx context.Context, docs []any, opts ...*options.InsertManyOptions) (*InsertManyResult, error) {
	sDocs := interfaceToSliceInterface(docs)
	if sDocs == nil {
		return nil, errors.New("must be valid slice to insert")
	}

	return c.collection.InsertMany(ctx, sDocs, opts...)
}
func (c *Collection) InsertOne(ctx context.Context, doc any, opts ...*options.InsertOneOptions) (*InsertOneResult, error) {
	return c.collection.InsertOne(ctx, doc, opts...)
}
func (c *Collection) ReplaceOne(ctx context.Context, query, doc any, opts ...*options.ReplaceOptions) (*UpdateResult, error) {
	return c.collection.ReplaceOne(ctx, query, doc, opts...)
}

func (c *Collection) UpdateAll(ctx context.Context, query, update any, opts ...*options.UpdateOptions) (*UpdateResult, error) {
	return c.collection.UpdateMany(ctx, query, update, opts...)
}

func (c *Collection) UpdateMany(ctx context.Context, query, update any, opts ...*options.UpdateOptions) (*UpdateResult, error) {
	return c.collection.UpdateMany(ctx, query, update, opts...)
}

func (c *Collection) UpdateOne(ctx context.Context, query, update any, opts ...*options.UpdateOptions) (*UpdateResult, error) {
	return c.collection.UpdateOne(ctx, query, update, opts...)
}

func (c *Collection) UpdateId(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*UpdateResult, error) {
	return c.collection.UpdateOne(ctx, bson.M{"_id": id}, update, opts...)
}

func (c *Collection) Remove(ctx context.Context, q interface{}) error {
	_, err := c.collection.DeleteOne(ctx, q)
	return err
}

func (c *Collection) RemoveId(ctx context.Context, id interface{}) error {
	res, err := c.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if res != nil && res.DeletedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	if err != nil {
		return err
	}
	return err
}

func (c *Collection) RemoveAll(ctx context.Context, q interface{}) (*ChangeInfo, error) {
	res, err := c.collection.DeleteMany(ctx, q)
	if err != nil {
		return nil, err
	}

	return &ChangeInfo{Removed: int(res.DeletedCount)}, nil
}

// ensureIndex create multiple indexes on the collection and returns the names of
// Example：indexes = []string{"idx1", "-idx2", "idx3,idx4"}
// Three indexes will be created, index idx1 with ascending order, index idx2 with descending order, idex3 and idex4 are Compound ascending sort index
// Reference: https://docs.mongodb.com/manual/reference/command/createIndexes/
func (c *Collection) ensureIndex(ctx context.Context, indexes []IndexModel) error {
	var indexModels []mongo.IndexModel
	for _, idx := range indexes {
		var model mongo.IndexModel
		var keysDoc bson.D

		for _, field := range idx.Key {
			key, n := SplitSortField(field)

			keysDoc = append(keysDoc, bson.E{Key: key, Value: n})
		}
		model = mongo.IndexModel{
			Keys:    keysDoc,
			Options: idx.IndexOptions,
		}

		indexModels = append(indexModels, model)
	}

	if len(indexModels) == 0 {
		return nil
	}

	res, err := c.collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil || len(res) == 0 {
		fmt.Println("<MongoDB.C>: ", c.collection.Name(), " Index: ", indexes, " error: ", err, "res: ", res)
		return err
	}
	return nil
}

// CreateIndexes creates multiple indexes in collection
// If the Key in opts.IndexModel is []string{"name"}, means create index: name
// If the Key in opts.IndexModel is []string{"name","-age"} means create Compound indexes: name and -age
func (c *Collection) CreateIndexes(ctx context.Context, indexes []IndexModel) error {
	return c.ensureIndex(ctx, indexes)
}

// CreateOneIndex creates one index
// If the Key in opts.IndexModel is []string{"name"}, means create index name
// If the Key in opts.IndexModel is []string{"name","-age"} means create Compound index: name and -age
func (c *Collection) CreateOneIndex(ctx context.Context, index IndexModel) error {
	return c.ensureIndex(ctx, []IndexModel{index})
}

// DropAllIndexes drop all indexes on the collection except the index on the _id field
// if there is only _id field index on the collection, the function call will report an error
func (c *Collection) DropAllIndexes(ctx context.Context) (err error) {
	_, err = c.collection.Indexes().DropAll(ctx)
	return err
}

// DropIndex drop indexes in collection, indexes that be dropped should be in line with inputting indexes
// The indexes is []string{"name"} means drop index: name
// The indexes is []string{"name","-age"} means drop Compound indexes: name and -age
func (c *Collection) DropIndex(ctx context.Context, indexes []string) error {
	_, err := c.collection.Indexes().DropOne(ctx, generateDroppedIndex(indexes))
	if err != nil {
		return err
	}
	return err
}

// generate indexes that store in mongo which may consist more than one index(like []string{"index1","index2"} is stored as "index1_1_index2_1")
func generateDroppedIndex(index []string) string {
	var res string
	for _, e := range index {
		key, sort := SplitSortField(e)
		n := key + "_" + fmt.Sprint(sort)
		if len(res) == 0 {
			res = n
		} else {
			res += "_" + n
		}
	}
	return res
}

// interfaceToSliceInterface convert interface to slice interface
func interfaceToSliceInterface(docs interface{}) []interface{} {
	if reflect.Slice != reflect.TypeOf(docs).Kind() {
		return nil
	}
	s := reflect.ValueOf(docs)
	if s.Len() == 0 {
		return nil
	}
	var sDocs []interface{}
	for i := 0; i < s.Len(); i++ {
		sDocs = append(sDocs, s.Index(i).Interface())
	}
	return sDocs
}
