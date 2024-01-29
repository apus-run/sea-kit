package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCursor(t *testing.T) {
	type QueryTestItem struct {
		Id   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
		Age  int                `bson:"age"`

		Instock []struct {
			Warehouse string `bson:"warehouse"`
			Qty       int    `bson:"qty"`
		} `bson:"instock"`
	}

	ast := require.New(t)
	cli := initClient()
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	id4 := primitive.NewObjectID()
	docs := []interface{}{
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, err := cli.InsertMany(context.Background(), docs)
	ast.NoError(err)

	var res QueryTestItem

	// if query has 1 record，cursor can run Next one time， Next time return false
	filter1 := bson.M{
		"name": "Alice",
	}

	cursor := cli.Find(context.Background(), filter1)
	ast.NoError(cursor.Err())

	val := cursor.Next(&res)
	ast.Equal(true, val)
	ast.Equal(id2, res.Id)

	val = cursor.Next(&res)
	ast.Equal(false, val)

	cursor.Close()

}
