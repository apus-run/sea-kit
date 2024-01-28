package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCollection_Insert(t *testing.T) {
	ast := require.New(t)
	ctx := context.Background()

	cli := initClient()

	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	cli.CreateIndexes(ctx, []IndexModel{
		{Key: []string{"id2", "id3"}},
		{Key: []string{"id4", "-id5"}},
	})

	var err error
	doc := bson.M{"_id": primitive.NewObjectID(), "name": "Alice"}

	res, err := cli.InsertOne(context.Background(), doc)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(doc["_id"], res.InsertedID)

	res, err = cli.InsertOne(context.Background(), doc)
	ast.Equal(true, IsDup(err))
	ast.Empty(res)
}

func TestCollection_CreateIndexes(t *testing.T) {
	ast := require.New(t)
	cli := initClient()
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	var expireS int32 = 100
	unique := []string{"id1"}
	indexOpts := options.Index()
	indexOpts.SetUnique(true).SetExpireAfterSeconds(expireS)
	ast.NoError(cli.CreateOneIndex(context.Background(), IndexModel{Key: unique, IndexOptions: indexOpts}))

	ast.NoError(cli.CreateIndexes(context.Background(), []IndexModel{{Key: []string{"id2", "id3"}},
		{Key: []string{"id4", "-id5"}}}))
	// same index，error
	ast.Error(cli.CreateOneIndex(context.Background(), IndexModel{Key: unique}))

	// check if unique indexs is working
	var err error
	doc := bson.M{
		"id1": 1,
	}

	_, err = cli.InsertOne(context.Background(), doc)
	ast.NoError(err)
	_, err = cli.InsertOne(context.Background(), doc)
	ast.Equal(true, IsDup(err))
}

func TestCollection_DropAllIndexes(t *testing.T) {
	ast := require.New(t)

	cli := initClient()
	defer cli.DropCollection(context.Background())

	var err error
	err = cli.DropAllIndexes(context.Background())
	ast.Error(err)
}

func TestCollection_DropIndex(t *testing.T) {
	ast := require.New(t)

	cli := initClient()
	defer cli.DropCollection(context.Background())

	indexOpts := options.Index()
	indexOpts.SetUnique(true)
	cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"index1"}, IndexOptions: indexOpts}})

	// same index，error
	ast.Error(cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"index1"}}}))

	err := cli.DropIndex(context.Background(), []string{"index1"})
	ast.NoError(err)
	ast.NoError(cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"index1"}}}))

	indexOpts = options.Index()
	indexOpts.SetUnique(true)
	cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"-index1"}, IndexOptions: indexOpts}})
	// same index，error
	ast.Error(cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"-index1"}}}))

	err = cli.DropIndex(context.Background(), []string{"-index1"})
	ast.NoError(err)
	ast.NoError(cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"-index1"}}}))

	err = cli.DropIndex(context.Background(), []string{""})
	ast.Error(err)

	err = cli.DropIndex(context.Background(), []string{"index2"})
	ast.Error(err)

	indexOpts = options.Index()
	indexOpts.SetUnique(true)
	cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"index2", "-index1"}, IndexOptions: indexOpts}})
	ast.Error(cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"index2", "-index1"}}}))
	err = cli.DropIndex(context.Background(), []string{"index2", "-index1"})
	ast.NoError(err)
	ast.NoError(cli.ensureIndex(context.Background(), []IndexModel{{Key: []string{"index2", "-index1"}}}))

	err = cli.DropIndex(context.Background(), []string{"-"})
	ast.Error(err)
}

func TestCollection_InsertMany(t *testing.T) {
	type UserInfo struct {
		Id     primitive.ObjectID `bson:"_id"`
		Name   string             `bson:"name"`
		Age    uint16             `bson:"age"`
		Weight uint32             `bson:"weight"`
	}

	ctx := context.Background()
	ast := require.New(t)
	cli := initClient()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	var err error
	newDocs := []any{UserInfo{Id: NewObjectID(), Name: "Alice", Age: 10}, UserInfo{Id: NewObjectID(), Name: "Lucas", Age: 11}}

	res, err := cli.InsertMany(ctx, newDocs)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(2, len(res.InsertedIDs))

	newPDocs := []any{UserInfo{Id: NewObjectID(), Name: "Alice3", Age: 10}, UserInfo{Id: NewObjectID(), Name: "Lucas3", Age: 11}}
	res, err = cli.InsertMany(ctx, newPDocs)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(2, len(res.InsertedIDs))

	docs4 := []UserInfo{}
	res, err = cli.InsertMany(context.Background(), []interface{}{docs4})
	ast.Error(err)
	ast.Empty(res)
}

func TestSliceInsert(t *testing.T) {
	type UserInfo struct {
		Id     primitive.ObjectID `bson:"_id"`
		Name   string             `bson:"name"`
		Age    uint16             `bson:"age"`
		Weight uint32             `bson:"weight"`
	}
	newDocs := []UserInfo{{Name: "Alice", Age: 10}, {Name: "Lucas", Age: 11}}
	di := interface{}(newDocs)
	dis := interfaceToSliceInterface(di)
	ast := require.New(t)
	ast.Len(dis, 2)

	newDocs_1 := []interface{}{UserInfo{Name: "Alice", Age: 10}, UserInfo{Name: "Lucas", Age: 11}}
	di = interface{}(newDocs_1)
	dis = interfaceToSliceInterface(di)
	ast.Len(dis, 2)

	newDocs_2 := UserInfo{Name: "Alice", Age: 10}
	di = interface{}(newDocs_2)
	dis = interfaceToSliceInterface(di)
	ast.Nil(dis)

	newDocs_3 := []UserInfo{}
	di = interface{}(newDocs_3)
	dis = interfaceToSliceInterface(di)
	ast = require.New(t)
	ast.Nil(dis)
}

func TestCollection_Update(t *testing.T) {
	ast := require.New(t)
	cli := initClient()
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Lucas"}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// update already exist record
	filter1 := bson.M{
		"name": "Alice",
	}
	update1 := bson.M{
		"$set": bson.M{
			"name": "Alice1",
			"age":  18,
		},
	}

	res, err := cli.UpdateOne(context.Background(), filter1, update1)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.NoError(err)

	// bloom_filter is nil or wrong BSON Document format
	update3 := bson.M{
		"$set": bson.M{
			"name": "Geek",
			"age":  21,
		},
	}
	res, err = cli.UpdateOne(context.Background(), nil, update3)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.Error(err)

	opts := options.Update().SetBypassDocumentValidation(false)
	res, err = cli.UpdateOne(context.Background(), 1, update3, opts)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.Error(err)

	// update is nil or wrong BSON Document format
	filter4 := bson.M{
		"name": "Geek",
	}
	res, err = cli.UpdateOne(context.Background(), filter4, nil)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.Error(err)

	res, err = cli.UpdateOne(context.Background(), filter4, 1)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.Error(err)
}

func TestCollection_UpdateId(t *testing.T) {
	ast := require.New(t)
	cli := initClient()
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Lucas"}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// update already exist record
	update1 := bson.M{
		"$set": bson.M{
			"name": "Alice1",
			"age":  18,
		},
	}

	opts := options.Update().SetBypassDocumentValidation(false)
	res, err := cli.UpdateId(context.Background(), id1, update1, opts)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.NoError(err)

	// id is nil or not exist
	update3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	res, err = cli.UpdateId(context.Background(), nil, update3)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.Error(err)

	res, err = cli.UpdateId(context.Background(), 1, update3)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.Error(err)

	res, err = cli.UpdateId(context.Background(), "not_exist_id", nil)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.Error(err)

	res, err = cli.UpdateId(context.Background(), "not_exist_id", 1)
	if res != nil && res.MatchedCount == 0 {
		err = mongo.ErrNoDocuments
	}
	ast.Error(err)
}

func TestCollection_UpdateAll(t *testing.T) {
	ast := require.New(t)
	cli := initClient()
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		bson.D{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// update already exist record
	filter1 := bson.M{
		"name": "Alice",
	}
	update1 := bson.M{
		"$set": bson.M{
			"age": 33,
		},
	}

	opts := options.Update().SetBypassDocumentValidation(false)
	res, err := cli.UpdateAll(context.Background(), filter1, update1, opts)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(2), res.MatchedCount)
	ast.Equal(int64(2), res.ModifiedCount)
	ast.Equal(int64(0), res.UpsertedCount)
	ast.Equal(nil, res.UpsertedID)

	// if record is not exist，err is nil， MatchedCount in res is 0
	filter2 := bson.M{
		"name": "Lily",
	}
	update2 := bson.M{
		"$set": bson.M{
			"age": 22,
		},
	}
	res, err = cli.UpdateAll(context.Background(), filter2, update2)
	ast.Nil(err)
	ast.NotNil(res)
	ast.Equal(int64(0), res.MatchedCount)

	// bloom_filter is nil or wrong BSON Document format
	update3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	res, err = cli.UpdateAll(context.Background(), nil, update3)
	ast.Error(err)
	ast.Nil(res)

	res, err = cli.UpdateAll(context.Background(), 1, update3)
	ast.Error(err)
	ast.Nil(res)

	// update is nil or wrong BSON Document format
	filter4 := bson.M{
		"name": "Geek",
	}
	res, err = cli.UpdateAll(context.Background(), filter4, nil)
	ast.Error(err)
	ast.Nil(res)

	res, err = cli.UpdateAll(context.Background(), filter4, 1)
	ast.Error(err)
	ast.Nil(res)
}

func TestCollection_Remove(t *testing.T) {
	ast := require.New(t)
	cli := initClient()
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	id1 := primitive.NewObjectID().Hex()
	id2 := primitive.NewObjectID().Hex()
	id3 := primitive.NewObjectID().Hex()
	id4 := primitive.NewObjectID().Hex()
	id5 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		bson.D{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
		bson.D{{Key: "_id", Value: id4}, {Key: "name", Value: "Joe"}, {Key: "age", Value: 20}},
		bson.D{{Key: "_id", Value: id5}, {Key: "name", Value: "Ethan"}, {Key: "age", Value: 1}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// remove id
	err = cli.RemoveId(context.Background(), "")
	ast.Error(err)
	err = cli.RemoveId(context.Background(), "not-exists-id")
	ast.True(IsErrNoDocuments(err))
	ast.NoError(cli.RemoveId(context.Background(), id4))
	ast.NoError(cli.RemoveId(context.Background(), id5))

	// delete record: name = "Alice" , after that, expect one name = "Alice" record
	filter1 := bson.M{
		"name": "Alice",
	}

	err = cli.Remove(context.Background(), filter1)
	ast.NoError(err)

	err = cli.Find(context.Background(), filter1).Err()
	ast.NoError(err)

	// delete not match  record , report err
	filter2 := bson.M{
		"name": "Lily",
	}
	err = cli.Remove(context.Background(), filter2)
	ast.Equal(err, mongo.ErrNoDocuments)

	// bloom_filter is bson.M{}，delete one document
	filter3 := bson.M{}
	err = cli.Find(context.Background(), filter3).Err()
	ast.NoError(err)

	err = cli.Remove(context.Background(), filter3)
	ast.NoError(err)

	err = cli.Find(context.Background(), filter3).Err()
	ast.NoError(err)

	// bloom_filter is nil or wrong BSON Document format
	err = cli.Remove(context.Background(), nil)
	ast.Error(err)

	err = cli.Remove(context.Background(), 1)
	ast.Error(err)
}

func TestCollection_RemoveAll(t *testing.T) {
	ast := require.New(t)
	cli := initClient()
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	id4 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		bson.D{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
		bson.D{{Key: "_id", Value: id4}, {Key: "name", Value: "Rocket"}, {Key: "age", Value: 23}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// delete record: name = "Alice" ,after that, expect - record : name = "Alice"
	filter1 := bson.M{
		"name": "Alice",
	}

	res, err := cli.RemoveAll(context.Background(), filter1)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal(2, res.Removed)

	err = cli.Find(context.Background(), filter1).Err()
	ast.NoError(err)

	// delete with not match bloom_filter， DeletedCount in res is 0
	filter2 := bson.M{
		"name": "Lily",
	}
	res, err = cli.RemoveAll(context.Background(), filter2)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal(0, res.Removed)

	// bloom_filter is bson.M{}，delete all docs
	filter3 := bson.M{}
	err = cli.Find(context.Background(), filter3).Err()
	ast.NoError(err)

	res, err = cli.RemoveAll(context.Background(), filter3)
	ast.NoError(err)
	ast.NotNil(res)

	err = cli.Find(context.Background(), filter3).Err()
	ast.NoError(err)

	// bloom_filter is nil or wrong BSON Document format
	res, err = cli.RemoveAll(context.Background(), nil)
	ast.Error(err)
	ast.Nil(res)

	res, err = cli.RemoveAll(context.Background(), 1)
	ast.Error(err)
	ast.Nil(res)
}

func TestCollection_ReplaceOne(t *testing.T) {
	type UserInfo struct {
		Id     primitive.ObjectID `bson:"_id"`
		Name   string             `bson:"name"`
		Age    uint16             `bson:"age"`
		Weight uint32             `bson:"weight"`
	}
	ast := require.New(t)
	cli := initClient()
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	id := primitive.NewObjectID()
	ui := UserInfo{Id: id, Name: "Lucas", Age: 17}
	_, err := cli.InsertOne(context.Background(), ui)
	ast.NoError(err)
	ui.Id = id
	ui.Age = 27
	_, err = cli.ReplaceOne(context.Background(), bson.M{"_id": id}, &ui)

	ast.NoError(err)

}
