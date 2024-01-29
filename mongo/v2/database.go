package mongo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is a handle to a MongoDB database
type Database struct {
	mu       sync.Mutex
	database *mongo.Database

	registry *bsoncodec.Registry

	processor processor
	logMode   bool
}

func (d *Database) SetLogMode(logMode bool) {
	d.logMode = logMode
}

func (d *Database) Client() *Client {
	d.mu.Lock()
	defer d.mu.Lock()

	cc := d.database.Client()
	if cc == nil {
		return nil
	}

	return &Client{
		client:   cc,
		registry: d.registry,
	}
}
func (d *Database) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *SingleResult {
	option := options.RunCmd()
	if len(opts) > 0 && opts != nil {
		option = opts[0]
	}
	return d.database.RunCommand(ctx, runCommand, option)
}
func (d *Database) RunCommandCursor(ctx context.Context, cmd interface{}) *Cursor {
	cur, err := d.database.RunCommandCursor(ctx, cmd)
	return &Cursor{
		ctx:    ctx,
		cursor: cur,
		err:    err,
	}
}

func (d *Database) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (
	cur *mongo.Cursor, err error) {
	err = d.processor(func(cmd *cmd) error {
		cur, err = d.database.ListCollections(ctx, filter, opts...)
		logCmd(d.logMode, cmd, "ListCollections", cur, filter)
		return err
	})
	return
}

// Collection gets collection from database
func (d *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	if d.database == nil {
		return nil
	}
	coll := d.database.Collection(name, opts...)
	if coll == nil {
		return nil
	}
	return &Collection{
		collection: coll,
		registry:   d.registry,

		processor: d.processor,
		logMode:   d.logMode,
	}
}

func (d *Database) CreateCollection(ctx context.Context, name string, opts *options.CollectionOptions) (*Collection, error) {
	if err := d.database.CreateCollection(ctx, name); err != nil {
		return nil, err
	}
	cp := d.database.Collection(name, opts)
	return &Collection{
		collection: cp,
		registry:   d.registry,
	}, nil
}

// GetDatabaseName returns the name of database
func (d *Database) GetDatabaseName() string {
	return d.database.Name()
}

func (d *Database) GetReadConcern() *readconcern.ReadConcern { return d.database.ReadConcern() }
func (d *Database) GetReadPreference() *readpref.ReadPref    { return d.database.ReadPreference() }

// DropDatabase drops database
func (d *Database) DropDatabase(ctx context.Context) error {
	return d.processor(func(c *cmd) error {
		logCmd(d.logMode, c, "Drop", nil)
		return d.database.Drop(ctx)
	})
}

func (d *Database) GetDatabase() *mongo.Database {
	return d.database
}

func (d *Database) WriteConcern() (res *writeconcern.WriteConcern) {
	_ = d.processor(func(c *cmd) error {
		res = d.database.WriteConcern()
		logCmd(d.logMode, c, "WriteConcern", res)
		return nil
	})
	return
}
