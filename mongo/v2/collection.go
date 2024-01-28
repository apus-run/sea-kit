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

type processor func(fn processFn) error
type processFn func(*cmd) error

type cmd struct {
	name     string
	req      []interface{}
	res      interface{}
	dbName   string
	collName string
}

func logCmd(logMode bool, c *cmd, name string, res any, req ...any) {
	// 只有开启log模式才会记录req、res
	if logMode {
		c.name = name
		c.req = append(c.req, req...)
		switch res := res.(type) {
		case *mongo.SingleResult:
			val, _ := res.Raw()
			c.res = val
		default:
			c.res = res
		}
	}
}

type Collection struct {
	collection *mongo.Collection

	registry *bsoncodec.Registry

	processor processor
	logMode   bool
}

func (c *Collection) cmd(cmd *cmd) *cmd {
	cmd.dbName = c.collection.Database().Name()
	cmd.collName = c.collection.Name()
	return cmd
}

func (c *Collection) SetLogMode(logMode bool) {
	c.logMode = logMode
}

func (c *Collection) Collection() *mongo.Collection {
	return c.collection
}

// DropCollection drops collection
// it's safe even collection is not exists
func (c *Collection) DropCollection(ctx context.Context) error {
	return c.processor(func(cmd *cmd) error {
		logCmd(c.logMode, c.cmd(cmd), "DropCollection", nil)
		return c.collection.Drop(ctx)
	})
}

// CloneCollection creates a copy of the Collection
func (c *Collection) CloneCollection() (coll *mongo.Collection, err error) {
	err = c.processor(func(cmd *cmd) error {
		coll, err = c.collection.Clone()
		logCmd(c.logMode, c.cmd(cmd), "CloneCollection", coll)

		return err
	})

	return coll, err

}

func (c *Collection) GetCollectionName() string {
	return c.collection.Name()
}

func (c *Collection) Aggregate(ctx context.Context, pipe any, opts ...*options.AggregateOptions) *Cursor {
	var cur *mongo.Cursor
	var err error
	err = c.processor(func(cmd *cmd) error {
		cur, err = c.collection.Aggregate(ctx, pipe, opts...)
		logCmd(c.logMode, cmd, "Aggregate", cur, pipe, opts)

		return err
	})

	return &Cursor{
		ctx:    ctx,
		cursor: cur,
		err:    err,
	}
}

func (c *Collection) Find(ctx context.Context, query any, opts ...*options.FindOptions) *Cursor {
	var cur *mongo.Cursor
	var err error

	err = c.processor(func(cmd *cmd) error {
		cur, err = c.collection.Find(ctx, query, opts...)
		logCmd(c.logMode, cmd, "Find", cur, query, opts)
		return err
	})

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

func (c *Collection) FindOne(ctx context.Context, query any, opts ...*options.FindOneOptions) (sr *SingleResult) {
	_ = c.processor(func(cmd *cmd) error {
		sr = c.collection.FindOne(ctx, query, opts...)
		logCmd(c.logMode, c.cmd(cmd), "FindOne", sr, query, opts)

		return sr.Err()
	})

	return
}

func (c *Collection) FindOneAndDelete(ctx context.Context, query any, opts ...*options.FindOneAndDeleteOptions) (sr *SingleResult) {
	_ = c.processor(func(cmd *cmd) error {
		sr = c.collection.FindOneAndDelete(ctx, query, opts...)
		logCmd(c.logMode, c.cmd(cmd), "FindOneAndDelete", sr, query, opts)

		return sr.Err()
	})
	return
}

func (c *Collection) FindOneAndReplace(ctx context.Context, query, doc any, opts ...*options.FindOneAndReplaceOptions) (sr *SingleResult) {
	_ = c.processor(func(cmd *cmd) error {
		sr = c.collection.FindOneAndReplace(ctx, query, doc, opts...)
		logCmd(c.logMode, c.cmd(cmd), "FindOneAndReplace", sr, query, doc, opts)

		return sr.Err()
	})

	return
}

func (c *Collection) FindOneAndUpdate(ctx context.Context, query, update any, opts ...*options.FindOneAndUpdateOptions) (sr *SingleResult) {
	_ = c.processor(func(cmd *cmd) error {
		sr = c.collection.FindOneAndUpdate(ctx, query, update, opts...)
		logCmd(c.logMode, c.cmd(cmd), "FindOneAndUpdate", sr, query, update, opts)

		return sr.Err()
	})

	return
}

func (c *Collection) Indexes() mongo.IndexView { return c.collection.Indexes() }

func (c *Collection) InsertMany(ctx context.Context, docs []any, opts ...*options.InsertManyOptions) (imr *InsertManyResult, err error) {

	sDocs := interfaceToSliceInterface(docs)
	if sDocs == nil {
		return nil, errors.New("must be valid slice to insert")
	}

	err = c.processor(func(cmd *cmd) error {
		imr, err = c.collection.InsertMany(ctx, sDocs, opts...)
		logCmd(c.logMode, c.cmd(cmd), "InsertMany", imr, docs, opts)

		return err
	})

	return
}
func (c *Collection) InsertOne(ctx context.Context, doc any, opts ...*options.InsertOneOptions) (ior *InsertOneResult, err error) {
	err = c.processor(func(cmd *cmd) error {
		ior, err = c.collection.InsertOne(ctx, doc, opts...)
		logCmd(c.logMode, c.cmd(cmd), "InsertOne", ior, doc, opts)
		return err
	})
	return
}
func (c *Collection) ReplaceOne(ctx context.Context, query, doc any, opts ...*options.ReplaceOptions) (ur *UpdateResult, err error) {
	err = c.processor(func(cmd *cmd) error {
		ur, err = c.collection.ReplaceOne(ctx, query, doc, opts...)
		logCmd(c.logMode, c.cmd(cmd), "ReplaceOne", ur, query, doc, opts)
		return err
	})

	return
}

func (c *Collection) UpdateAll(ctx context.Context, query, update any, opts ...*options.UpdateOptions) (ur *UpdateResult, err error) {
	err = c.processor(func(cmd *cmd) error {
		ur, err = c.collection.UpdateMany(ctx, query, update, opts...)
		logCmd(c.logMode, c.cmd(cmd), "UpdateAll", ur, query, update, opts)
		return err
	})
	return
}

func (c *Collection) UpdateMany(ctx context.Context, query, update any, opts ...*options.UpdateOptions) (ur *UpdateResult, err error) {
	err = c.processor(func(cmd *cmd) error {
		ur, err = c.collection.UpdateMany(ctx, query, update, opts...)
		logCmd(c.logMode, c.cmd(cmd), "UpdateMany", ur, query, update, opts)
		return err
	})
	return
}

func (c *Collection) UpdateOne(ctx context.Context, query, update any, opts ...*options.UpdateOptions) (ur *UpdateResult, err error) {
	err = c.processor(func(cmd *cmd) error {
		ur, err = c.collection.UpdateOne(ctx, query, update, opts...)
		logCmd(c.logMode, c.cmd(cmd), "UpdateOne", ur, query, update, opts)
		return err
	})

	return
}

func (c *Collection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	_ = c.processor(func(cmd *cmd) error {
		res, err = c.collection.UpdateByID(ctx, id, update, opts...)
		logCmd(c.logMode, c.cmd(cmd), "UpdateByID", res, id, update)
		return err
	})
	return
}

func (c *Collection) UpdateId(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (ur *UpdateResult, err error) {
	err = c.processor(func(cmd *cmd) error {
		ur, err = c.collection.UpdateOne(ctx, bson.M{"_id": id}, update, opts...)
		logCmd(c.logMode, c.cmd(cmd), "UpdateId", ur, id, update, opts)
		return err
	})
	return
}

func (c *Collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (
	res *DeleteResult, err error) {

	err = c.processor(func(cmd *cmd) error {
		res, err = c.collection.DeleteMany(ctx, filter, opts...)
		logCmd(c.logMode, c.cmd(cmd), "DeleteMany", res, filter)
		return err
	})
	return
}

func (c *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (dr *DeleteResult, err error) {
	err = c.processor(func(cmd *cmd) error {
		dr, err = c.collection.DeleteOne(ctx, filter, opts...)
		logCmd(c.logMode, c.cmd(cmd), "DeleteOne", dr, filter)
		return err
	})
	return
}

func (c *Collection) Remove(ctx context.Context, q interface{}) (err error) {
	var res *DeleteResult
	err = c.processor(func(cmd *cmd) error {
		res, err = c.collection.DeleteMany(ctx, q)
		logCmd(c.logMode, c.cmd(cmd), "Remove", res, q)
		return err
	})
	return
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
	var err error
	var res *DeleteResult
	err = c.processor(func(cmd *cmd) error {
		res, err = c.collection.DeleteMany(ctx, q)
		logCmd(c.logMode, c.cmd(cmd), "RemoveAll", res)
		return err
	})

	return &ChangeInfo{Removed: int(res.DeletedCount)}, err
}

func (c *Collection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) (res []interface{}, err error) {
	err = c.processor(func(cmd *cmd) error {
		res, err = c.collection.Distinct(ctx, fieldName, filter, opts...)
		logCmd(c.logMode, c.cmd(cmd), "Distinct", res, fieldName, filter)
		return err
	})
	return
}

func (c *Collection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (res int64, err error) {
	err = c.processor(func(cmd *cmd) error {
		res, err = c.collection.EstimatedDocumentCount(ctx, opts...)
		logCmd(c.logMode, c.cmd(cmd), "EstimatedDocumentCount", res)
		return err
	})
	return
}

func (c *Collection) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (res *mongo.ChangeStream, err error) {
	_ = c.processor(func(cmd *cmd) error {
		res, err = c.collection.Watch(ctx, pipeline, opts...)
		logCmd(c.logMode, c.cmd(cmd), "Watch", res, pipeline)
		return err
	})
	return
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
