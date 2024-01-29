// https://github.com/megaease/easeprobe/blob/main/probe/client/conf/conf.go
// https://github.com/megaease/easeprobe/blob/main/probe/client/client.go

package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// DriverType is the client driver
type DriverType int

const (
	Unknown DriverType = iota
	MySQL
	Redis
	Mongo
	PostgreSQL
)

// DriverMap is the safemap of [driver, name]
var DriverMap = map[DriverType]string{
	MySQL:      "mysql",
	Redis:      "redis",
	Mongo:      "mongo",
	PostgreSQL: "postgres",
	Unknown:    "unknown",
}

// DriverTypeMap is the safemap of driver [name, driver]
var DriverTypeMap = ReverseMap(DriverMap)

// String convert the DriverType to string
func (d DriverType) String() string {
	if val, ok := DriverMap[d]; ok {
		return val
	}
	return DriverMap[Unknown]
}

// DriverType convert the string to DriverType
func (d *DriverType) DriverType(name string) DriverType {
	if val, ok := DriverTypeMap[name]; ok {
		*d = val
		return val
	}
	return Unknown
}

// MarshalYAML is marshal the provider type
func (d DriverType) MarshalYAML() (interface{}, error) {
	return EnumMarshalYaml(DriverMap, d, "Test")
}

// UnmarshalYAML is unmarshal the provider type
func (d *DriverType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return EnumUnmarshalYaml(unmarshal, DriverTypeMap, d, Unknown, "Test")
}

// MarshalJSON is marshal the provider
func (d DriverType) MarshalJSON() (b []byte, err error) {
	return EnumMarshalJSON(DriverMap, d, "Test")
}

// UnmarshalJSON is Unmarshal the provider type
func (d *DriverType) UnmarshalJSON(b []byte) (err error) {
	return EnumUnmarshalJSON(b, DriverTypeMap, d, Unknown, "Test")
}

func testMarshalUnmarshal(t *testing.T, str string, te DriverType, good bool,
	marshal func(in interface{}) ([]byte, error),
	unmarshal func(in []byte, out interface{}) (err error)) {

	var s DriverType
	err := unmarshal([]byte(str), &s)
	if good {
		assert.Nil(t, err)
		assert.Equal(t, te, s)
	} else {
		assert.Error(t, err)
		assert.Equal(t, Unknown, s)
	}

	buf, err := marshal(te)
	if good {
		assert.Nil(t, err)
		assert.Equal(t, str, string(buf))
	} else {
		assert.Error(t, err)
		assert.Nil(t, buf)
	}
}
func testYamlJSON(t *testing.T, str string, te DriverType, good bool) {
	testYaml(t, str+"\n", te, good)
	testJSON(t, `"`+str+`"`, te, good)
}
func testYaml(t *testing.T, str string, te DriverType, good bool) {
	testMarshalUnmarshal(t, str, te, good, yaml.Marshal, yaml.Unmarshal)
}
func testJSON(t *testing.T, str string, te DriverType, good bool) {
	testMarshalUnmarshal(t, str, te, good, json.Marshal, json.Unmarshal)
}

func testDriverType(t *testing.T, str string, driver DriverType) {
	var d DriverType
	d.DriverType(str)
	assert.Equal(t, driver, d)

	s := driver.String()
	assert.Equal(t, str, s)
}

func TestDriverType(t *testing.T) {
	testDriverType(t, "mysql", MySQL)
	testDriverType(t, "redis", Redis)
	testDriverType(t, "mongo", Mongo)
	testDriverType(t, "postgres", PostgreSQL)
	testDriverType(t, "unknown", Unknown)

	d := Unknown
	assert.Equal(t, MySQL, d.DriverType("mysql"))
	assert.Equal(t, Redis, d.DriverType("redis"))

	d = 10
	assert.Equal(t, "unknown", d.String())
	assert.Equal(t, Unknown, d.DriverType("bad"))

	testYamlJSON(t, "mysql", MySQL, true)
	testYamlJSON(t, "redis", Redis, true)
	testYamlJSON(t, "mongo", Mongo, true)
	testYamlJSON(t, "postgres", PostgreSQL, true)
	testYamlJSON(t, "unknown", Unknown, true)

	testJSON(t, "", 10, false)
	testJSON(t, `{"x":"y"}`, 10, false)
	testJSON(t, `"xyz"`, 10, false)
	testYaml(t, "- mysql::", 10, false)
	testYamlJSON(t, "bad", 10, false)
	testJSON(t, `{"x":"y"}`, 10, false)
	testYaml(t, "-bad::", 10, false)
}
