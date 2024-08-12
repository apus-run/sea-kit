package jsonutil

import (
	"encoding/json"
	"strconv"
)

// Bytes is a byte slice in a json-encoded struct.
// encoding/json assumes that []byte fields are hex-encoded.
// Bytes are not hex-encoded; they are treated the same as strings.
// This can avoid unnecessary allocations due to a round trip through strings.
type Bytes []byte

func (b *Bytes) UnmarshalText(text []byte) error {
	// Copy the contexts of text.
	*b = append(*b, text...)
	return nil
}

// Int forces the given value to be an int during the JSON unmarshalling,
// even if float was given.
type Int int

func (i *Int) UnmarshalJSON(b []byte) error {
	// Handle null case.
	if string(b) == "null" {
		// Set default value
		*i = 0
		// Return nil
		return nil
	}
	// Parse string as float.
	// This approach will cover both int and float cases.
	v, err := strconv.ParseFloat(string(b), 64)
	// Set the value
	*i = Int(v)
	// Return error
	return err
}

/*
String is a function that converts the given json string into map[string]any.
Should be used in a known environment, when values are expected to be correct.
Panics in case of failure.

Usage:

	jsonx.Map({"foo":1,"bar":2}) // map[string]any{"foo": 1, "bar": 2}
*/
func Map(jsonstring string) map[string]any {
	data := map[string]any{}
	err := json.Unmarshal([]byte(jsonstring), &data)
	if err != nil {
		panic(err)
	}
	return data
}

/*
String is a function that converts the given value to a JSON string.
Almost useless for typical use-cases, but useful as template function.

Usage:

	jsonx.String(map[any]any{"foo": 1, "bar": 2}) // {"foo":1,"bar":2}
*/
func String(value any) string {
	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(v)
}
