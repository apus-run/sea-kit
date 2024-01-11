package mongo

import (
	"bytes"
	"errors"

	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestPrepSort(t *testing.T) {

	tbl := []struct {
		inp []string
		out bson.D
	}{
		{nil, bson.D{}},
		{[]string{"f1", " f2", "-f3 ", "+f4"}, bson.D{{"f1", 1}, {"f2", 1}, {"f3", -1}, {"f4", 1}}},
		{[]string{"+f1", " -f2", "-f3", " f4 "}, bson.D{{"f1", 1}, {"f2", -1}, {"f3", -1}, {"f4", 1}}},
	}

	for i, tt := range tbl {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out := PrepSort(tt.inp...)
			assert.EqualValues(t, tt.out, out)
		})
	}
}

func TestPrepIndex(t *testing.T) {
	tbl := []struct {
		inp []string
		out mongo.IndexModel
	}{
		{nil, mongo.IndexModel{Keys: bson.D{}}},
		{[]string{"f1", " f2", "-f3 ", "+f4"}, mongo.IndexModel{Keys: bson.D{{"f1", 1}, {"f2", 1}, {"f3", -1}, {"f4", 1}}}},
		{[]string{"+f1", " -f2", "-f3", " f4 "}, mongo.IndexModel{Keys: bson.D{{"f1", 1}, {"f2", -1}, {"f3", -1}, {"f4", 1}}}},
	}

	for i, tt := range tbl {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out := PrepIndex(tt.inp...)
			assert.EqualValues(t, tt.out, out)
		})
	}
}

func TestBind(t *testing.T) {
	type request struct {
		Fields     []string `json:"fields" bson:"fields"`
		Filter     bson.M   `json:"bloom_filter" bson:"bloom_filter"`
		Psrc       string   `json:"psrc" bson:"psrc"`
		SubTotals  bool     `json:"subtotals" bson:"subtotals"`
		StatFilter bson.M   `json:"stat_filter" bson:"stat_filter"`
		Sort       bson.D   `json:"sort" bson:"sort"`
		Dry        bool     `json:"dry" bson:"dry"`
		Encrypted  bool     `json:"-" bson:"-"`
	}
	body := bytes.NewBufferString(`{"fields":["cusip","acc"], "bloom_filter":{"trade_dt":{"$gte":{"$date":"2020-08-17T00:00:00-04:00"}, "$lt":{"$date":"2020-08-21T23:59:59-04:00"}}}, "psrc":"DEMO", "sort":{"day":1, "trade_id":1}, "stat_filter":{}, "subtotals":false}`)

	res := request{}
	err := Bind(body, &res)
	require.NoError(t, err)
	t.Logf("%+v", res)

	assert.Equal(t, []string{"cusip", "acc"}, res.Fields) // nolint
	assert.Equal(t, bson.M{"trade_dt": bson.M{"$gte": primitive.DateTime(1597636800000), "$lt": primitive.DateTime(1598068799000)}}, res.Filter)
	assert.Equal(t, bson.D{{"day", int32(1)}, {"trade_id", int32(1)}}, res.Sort)

	assert.Equal(t, "DEMO", res.Psrc)

	body = bytes.NewBufferString(`{"fields":["cusip","acc"], "bloom_filter":{"trade_dt":{"$gte":{"$date":"2020-08-17T00:00:00-04:00"}, "$lt":{"$date":"2020-08-21T23:59:59-04:00"}}}, "page":{"num":0, "size":50}, "psrc":"DEMO", "sort":{"trade_id":1, "day":1}, "stat_filter":{}, "subtotals":false}`)

	res = request{}
	err = Bind(body, &res)
	require.NoError(t, err)
	t.Logf("%+v", res)

	assert.Equal(t, []string{"cusip", "acc"}, res.Fields)
	assert.Equal(t, bson.M{"trade_dt": bson.M{"$gte": primitive.DateTime(1597636800000), "$lt": primitive.DateTime(1598068799000)}}, res.Filter)
	assert.Equal(t, bson.D{{"trade_id", int32(1)}, {"day", int32(1)}}, res.Sort)

	assert.Equal(t, "DEMO", res.Psrc)

	body = bytes.NewBufferString(`{"fields":["cusip","acc"], "bloom_filter":{"trade_dt":{"$gte":{"$date":"2020-08-17T04:00:00Z"}, "$lt":{"$date":"2020-08-21T23:59:59-04:00"}}}, "page":{"num":0, "size":50}, "psrc":"DEMO", "sort":{"trade_id":1, "day":1}, "stat_filter":{}, "subtotals":false}`)

	res = request{}
	err = Bind(body, &res)
	require.NoError(t, err)
	t.Logf("%+v", res)
	assert.Equal(t, bson.M{"trade_dt": bson.M{"$gte": primitive.DateTime(1597636800000), "$lt": primitive.DateTime(1598068799000)}}, res.Filter)
}

func TestIsErrNoDocuments(t *testing.T) {
	ast := require.New(t)
	ast.False(IsErrNoDocuments(errors.New("dont match")))
	ast.True(IsErrNoDocuments(mongo.ErrNoDocuments))
}

func TestIsDup(t *testing.T) {
	ast := require.New(t)
	ast.False(IsDup(nil))
	ast.False(IsDup(errors.New("invaliderror")))
	ast.True(IsDup(errors.New("E11000")))
}
