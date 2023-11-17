package mongo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

func initClient() *Mongo {
	ctx := context.Background()
	cli, err := Open(ctx, WithMongoConfig(func(options *Config) {
		monitor := &event.CommandMonitor{
			// 每个命令（查询）执行之前
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				log.Println("查询执行之前", startedEvent.Command)
			},
			// 执行成功
			Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
				log.Println("执行成功")
			},
			// 执行失败
			Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
				log.Println("执行失败")
			},
		}
		options.Uri = "mongodb://root:example@localhost:27017/"
		options.DatabaseName = "test_db"
		options.CollectionName = "testopen"
		//connectTimeout := 30 * time.Second
		//maxConnIdleTime := 3 * time.Minute
		//minPoolSize := uint64(20)
		//maxPoolSize := uint64(300)
		//options.ClientOptions.ConnectTimeout = &connectTimeout
		//options.ClientOptions.MaxConnIdleTime = &maxConnIdleTime
		//options.ClientOptions.MinPoolSize = &minPoolSize
		//options.ClientOptions.MaxPoolSize = &maxPoolSize

		options.Monitor = monitor
	}))

	if err != nil {
		panic(err)
	}

	return cli
}

func TestClient_ServerVersion(t *testing.T) {
	cli := initClient()

	version := cli.ServerVersion()
	assert.NotEmpty(t, version)
	t.Logf("version: %v", version)
}

func TestNewClientOpen(t *testing.T) {
	ctx := context.TODO()
	monitor := &event.CommandMonitor{
		// 每个命令（查询）执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			t.Log("查询执行之前", startedEvent.Command)
		},
		// 执行成功
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
			t.Log("执行成功")
		},
		// 执行失败
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
			t.Log("执行失败")
		},
	}
	cli, err := Open(
		ctx,
		WithURI("mongodb://root:example@localhost:27017/"),
		WithDatabaseName("test_db"),
		WithCollectionName("testopen"),
		WithClientOptions(&options.ClientOptions{
			Monitor: monitor,
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, cli.GetDatabaseName(), "test_db")
	assert.Equal(t, cli.GetCollectionName(), "testopen")

	err = cli.Ping(5)
	assert.NoError(t, err)

	res, err := cli.InsertOne(
		ctx,
		Article{
			Id:      123,
			Title:   "我的标题",
			Content: "我的内容",
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	t.Logf("结果: %v", res)

	filter := D(KV("id", 123))
	var art Article
	err = cli.FindOne(ctx, filter).Decode(&art)
	if err == mongo.ErrNoDocuments {
		t.Logf("没有找到 id: %v", 123)
	}
	t.Logf("查询结果: %#v", art)

	cli.DropCollection(context.Background())

	// close Client
	cli.Close(ctx)
	_, err = cli.InsertOne(
		ctx,
		Article{
			Id:      123,
			Title:   "我的标题",
			Content: "我的内容",
		},
	)
	assert.EqualError(t, err, "client is disconnected")

	err = cli.Ping(5)
	assert.Error(t, err)
}

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient(ctx, WithURI("mongodb://root:example@localhost:27017/"))
	if err != nil {
		t.Fatalf("连接客户端失败: %v", err)
	}

	db := client.Database("test_db")
	coll := db.Collection("testopen")

	assert.Equal(t, coll.GetCollectionName(), "testopen")

	err = client.Ping(5)
	assert.NoError(t, err)

	res, err := coll.InsertOne(ctx, D(KV("count", 1)))
	assert.NoError(t, err)
	assert.NotNil(t, res)
	t.Logf("结果: %v", res)

	coll.DropCollection(ctx)

	// close Client
	client.Close(ctx)
	_, err = coll.InsertOne(ctx, bson.D{{Key: "x", Value: 1}})
	assert.EqualError(t, err, "client is disconnected")

	err = client.Ping(5)
	assert.Error(t, err)
}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
