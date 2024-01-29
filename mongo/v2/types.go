package mongo

import (
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InsertOneResult = mongo.InsertOneResult
type InsertManyResult = mongo.InsertManyResult
type UpdateResult = mongo.UpdateResult
type SingleResult = mongo.SingleResult
type DeleteResult = mongo.DeleteResult

type IndexModel struct {
	Key []string // Index key fields; prefix name with dash (-) for descending order
	*options.IndexOptions
}

// ChangeInfo represents the data returned by Update and Upsert
// documents. This type mirrors the mgo type.
type ChangeInfo struct {
	Updated    int // Number of existing documents updated
	Removed    int // Number of documents removed
	UpsertedId any // Upserted _id field, when not explicitly provided
}

// Change represents the options that you can pass to the
// findAndModify operation.
type Change struct {
	Update    any  // The update document
	Upsert    bool // Whether to insert in case the document isn't found
	Remove    bool // Whether to remove the document found rather than updating
	ReturnNew bool // Should the modified document be returned rather than the old one
}

type KeyValue struct {
	Key   string
	Value any
}

func KV(key string, value any) KeyValue {
	return KeyValue{Key: key, Value: value}
}

func M(key string, value any) bson.M {
	return bson.M{key: value}
}

func E(key string, value any) bson.E {
	return bson.E{Key: key, Value: value}
}

func A[T any](values ...T) bson.A {
	value := make(bson.A, 0, len(values))
	for _, v := range values {
		value = append(value, v)
	}
	return value
}

func D(bsonElements ...KeyValue) bson.D {
	value := make(bson.D, 0, len(bsonElements))
	for _, element := range bsonElements {
		value = append(value, bson.E{Key: element.Key, Value: element.Value})
	}
	return value
}

func ID(value any) bson.M {
	return M("_id", value)
}

func ToUpdate(vals map[string]any) bson.M {
	return vals
}

func ToFilter(vals map[string]any) bson.D {
	var res bson.D
	for k, v := range vals {
		res = append(res, bson.E{k, v})
	}
	return res
}

func Set(vals map[string]any) bson.M {
	return bson.M{"$set": bson.M(vals)}
}

func GetSort(keys []string) bson.D {
	if len(keys) == 0 {
		return nil
	}

	sort := bson.D{}

	for _, k := range keys {
		if strings.HasPrefix(k, "-") {
			sort = append(sort, bson.E{Key: k[1:], Value: -1})
		} else if strings.HasPrefix(k, "+") {
			sort = append(sort, bson.E{Key: k[1:], Value: 1})
		} else {
			sort = append(sort, bson.E{Key: k, Value: 1})
		}
	}

	return sort
}

// PrepSort prepares sort params for mongo driver and returns bson.D
// Input string provided as [+|-]field1,[+|-]field2,[+|-]field3...
// + means ascending, - means descending. Lack of + or - in the beginning of the field name means ascending sort.
func PrepSort(sort ...string) bson.D {
	res := bson.D{}
	for _, s := range sort {
		if s == "" {
			continue
		}
		s = strings.TrimSpace(s)
		switch s[0] {
		case '-':
			res = append(res, bson.E{Key: s[1:], Value: -1})
		case '+':
			res = append(res, bson.E{Key: s[1:], Value: 1})
		default:
			res = append(res, bson.E{Key: s, Value: 1})
		}
	}
	return res
}

// PrepIndex prepares index params for mongo driver and returns IndexModel
func PrepIndex(keys ...string) mongo.IndexModel {
	return mongo.IndexModel{Keys: PrepSort(keys...)}
}

// Bind request json body from io.Reader to bson record
func Bind(r io.Reader, v interface{}) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return bson.UnmarshalExtJSON(body, false, v)
}

func GetFindAndModifyReturn(returnNew bool) options.ReturnDocument {
	if returnNew {
		return options.After
	}
	return options.Before
}

// IsErrNoDocuments check if err is no documents
func IsErrNoDocuments(err error) bool {
	if err == mongo.ErrNoDocuments {
		return true
	}
	return false
}

// IsDup check if err is mongo E11000 (duplicate err)。
func IsDup(err error) bool {
	return err != nil && strings.Contains(err.Error(), "E11000")
}

// SplitSortField handle sort symbol: "+"/"-" in front of field
// if "+"， return sort as 1
// if "-"， return sort as -1
func SplitSortField(field string) (key string, sort int32) {
	sort = 1
	key = field

	if len(field) != 0 {
		switch field[0] {
		case '+':
			key = strings.TrimPrefix(field, "+")
			sort = 1
		case '-':
			key = strings.TrimPrefix(field, "-")
			sort = -1
		}
	}

	return key, sort
}

// Now return Millisecond current time
func Now() time.Time {
	return time.Unix(0, time.Now().UnixNano()/1e6*1e6)
}

// NewObjectID generates a new ObjectID.
// Watch out: the way it generates objectID is different from mgo
func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}
