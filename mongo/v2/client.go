package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client creates client to mongo
type Client struct {
	client *mongo.Client

	registry *bsoncodec.Registry

	processor processor
	logMode   bool
}

func defaultProcessor(processFn processFn) error {
	return processFn(&cmd{req: make([]interface{}, 0, 1)})
}

// NewClient creates MongoDB client
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := Apply(opts...)
	return client(ctx, cfg)
}

func client(ctx context.Context, cfg *Config) (*Client, error) {
	option := options.Client()
	option = cfg.ClientOptions
	option.ApplyURI(cfg.Uri)

	client, err := mongo.Connect(ctx, option)
	if err != nil {
		return nil, fmt.Errorf("can't connect to mongo: %w", err)
	}

	// half of default connect timeout
	pCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err = client.Ping(pCtx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("can't ping mongo: %w", err)
	}

	return &Client{
		client:    client,
		registry:  option.Registry,
		processor: defaultProcessor,
	}, nil
}

func (c *Client) SetLogMode(logMode bool) {
	c.logMode = logMode
}

func (c *Client) WrapProcessor(wrapFn func(processFn) processFn) {
	c.processor = func(fn processFn) error {
		return wrapFn(fn)(&cmd{req: make([]interface{}, 0, 1)})
	}
}

// Disconnect closes sockets to the topology referenced by this Client.
func (c *Client) Disconnect(ctx context.Context) error {
	return c.processor(func(cmd *cmd) error {
		logCmd(c.logMode, cmd, "Disconnect", nil)
		return c.client.Disconnect(ctx)
	})
}

// Close closes sockets to the topology referenced by this Client.
func (c *Client) Close(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	return err
}

// Ping confirm connection is alive
func (c *Client) Ping(timeout int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	return c.processor(func(cmd *cmd) error {
		logCmd(c.logMode, cmd, "Ping", nil)
		return c.client.Ping(ctx, readpref.Primary())
	})
}

// Database create connection to database
func (c *Client) Database(name string, dbOpts ...*options.DatabaseOptions) *Database {
	var db *mongo.Database
	c.processor(func(cmd *cmd) error {
		db = c.client.Database(name, dbOpts...)
		cmd.dbName = name
		logCmd(c.logMode, cmd, "Database", db, name)
		return nil
	})
	if db == nil {
		return nil
	}
	return &Database{
		database: db,
		registry: c.registry,

		processor: c.processor,
		logMode:   c.logMode,
	}
}

// ServerVersion get the version of mongoDB server, like 4.4.0
func (c *Client) ServerVersion() string {
	var buildInfo bson.Raw
	err := c.client.Database("admin").RunCommand(
		context.Background(),
		bson.D{{"buildInfo", 1}},
	).Decode(&buildInfo)
	if err != nil {
		fmt.Println("run command err", err)
		return ""
	}
	v, err := buildInfo.LookupErr("version")
	if err != nil {
		fmt.Println("look up err", err)
		return ""
	}
	return v.StringValue()
}

func (c *Client) Client() *mongo.Client { return c.client }

func (c *Client) StartSession(opts ...*options.SessionOptions) (ss Session, err error) {
	err = c.processor(func(cmd *cmd) error {
		ss, err = c.client.StartSession(opts...)
		logCmd(c.logMode, cmd, "StartSession", ss)
		return err
	})
	return &session{Session: ss, logMode: c.logMode, processor: c.processor}, nil
}

func (c *Client) UseSession(ctx context.Context, fn func(SessionContext) error) error {
	return c.processor(func(cmd *cmd) error {
		logCmd(c.logMode, cmd, "UseSession", nil)
		return c.client.UseSession(ctx, fn)
	})
}

func (c *Client) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(SessionContext) error) error {
	return c.processor(func(cmd *cmd) error {
		logCmd(c.logMode, cmd, "UseSessionWithOptions", nil)
		return c.client.UseSessionWithOptions(ctx, opts, fn)
	})
}

func (c *Client) ListDatabaseNames(ctx context.Context, filter any, opts ...*options.ListDatabasesOptions) (
	dbs []string, err error) {

	err = c.processor(func(cmd *cmd) error {
		dbs, err = c.client.ListDatabaseNames(ctx, filter, opts...)
		logCmd(c.logMode, cmd, "ListDatabaseNames", dbs, filter)
		return err
	})
	return
}

func (c *Client) ListDatabases(ctx context.Context, filter any, opts ...*options.ListDatabasesOptions) (
	dbr mongo.ListDatabasesResult, err error) {

	err = c.processor(func(cmd *cmd) error {
		dbr, err = c.client.ListDatabases(ctx, filter, opts...)
		logCmd(c.logMode, cmd, "ListDatabases", dbr, filter)
		return err
	})
	return
}

func WithSession(ctx context.Context, sess Session, fn func(SessionContext) error) error {
	return mongo.WithSession(ctx, sess, fn)
}
