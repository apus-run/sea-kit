package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_Client(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient(ctx, WithURI("mongodb://root:example@localhost:27017/"))
	if err != nil {
		t.Fatalf("连接客户端失败: %v", err)
	}

	assert.NoError(t, err)
	db := client.Database("test_db")
	assert.Nil(t, err)
	assert.Equal(t, "test_db", db.GetDatabaseName())
	coll := db.Collection("testopen")
	assert.Equal(t, "testopen", coll.GetCollectionName())
	db.Collection("testopen").DropCollection(context.Background())
	db.DropDatabase(context.Background())
}

func TestRunCommand(t *testing.T) {
	ast := require.New(t)
	ctx := context.Background()
	cli := initClient()

	opts := options.RunCmd().SetReadPreference(readpref.Primary())
	res := cli.RunCommand(ctx, bson.D{{"ping", 1}}, opts)
	ast.NoError(res.Err())

	defer cli.Close(ctx)
}
