package collection

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// AnyValue 类型转换结构定义
type AnyValue struct {
	Value any
	Error error
}

// Int 返回 int 数据
func (av AnyValue) Int() (int, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(int)
	if !ok {
		return 0, NewErrInvalidType("int", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsInt() (int, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	switch v := av.Value.(type) {
	case int:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		return int(res), err
	}
	return 0, NewErrInvalidType("int", av.Value)
}

// IntOrDefault 返回 int 数据，或者默认值
func (av AnyValue) IntOrDefault(def int) int {
	val, err := av.Int()
	if err != nil {
		return def
	}
	return val
}

// Uint 返回 uint 数据
func (av AnyValue) Uint() (uint, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(uint)
	if !ok {
		return 0, NewErrInvalidType("uint", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsUint() (uint, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	switch v := av.Value.(type) {
	case uint:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 64)
		return uint(res), err
	}
	return 0, NewErrInvalidType("uint", av.Value)
}

// UintOrDefault 返回 uint 数据，或者默认值
func (av AnyValue) UintOrDefault(def uint) uint {
	val, err := av.Uint()
	if err != nil {
		return def
	}
	return val
}

func (av AnyValue) Int8() (int8, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(int8)
	if !ok {
		return 0, NewErrInvalidType("int", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsInt8() (int8, error) {
	if av.Error != nil {
		return 0, av.Error
	}

	switch v := av.Value.(type) {
	case int8:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		return int8(res), err
	}
	return 0, NewErrInvalidType("int8", av.Value)
}

func (av AnyValue) Int8OrDefault(def int8) int8 {
	val, err := av.Int8()
	if err != nil {
		return def
	}
	return val
}

func (av AnyValue) Uint8() (uint8, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(uint8)
	if !ok {
		return 0, NewErrInvalidType("uint8", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsUint8() (uint8, error) {
	if av.Error != nil {
		return 0, av.Error
	}

	switch v := av.Value.(type) {
	case uint8:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 8)
		return uint8(res), err
	}
	return 0, NewErrInvalidType("uint8", av.Value)
}

func (av AnyValue) Uint8OrDefault(def uint8) uint8 {
	val, err := av.Uint8()
	if err != nil {
		return def
	}
	return val
}

func (av AnyValue) Int16() (int16, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(int16)
	if !ok {
		return 0, NewErrInvalidType("int16", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsInt16() (int16, error) {
	if av.Error != nil {
		return 0, av.Error
	}

	switch v := av.Value.(type) {
	case int16:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 16)
		return int16(res), err
	}
	return 0, NewErrInvalidType("int16", av.Value)
}

func (av AnyValue) Int16OrDefault(def int16) int16 {
	val, err := av.Int16()
	if err != nil {
		return def
	}
	return val
}

func (av AnyValue) Uint16() (uint16, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(uint16)
	if !ok {
		return 0, NewErrInvalidType("uint16", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsUint16() (uint16, error) {
	if av.Error != nil {
		return 0, av.Error
	}

	switch v := av.Value.(type) {
	case uint16:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 16)
		return uint16(res), err
	}
	return 0, NewErrInvalidType("uint16", av.Value)
}

func (av AnyValue) Uint16OrDefault(def uint16) uint16 {
	val, err := av.Uint16()
	if err != nil {
		return def
	}
	return val
}

// Int32 返回 int32 数据
func (av AnyValue) Int32() (int32, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(int32)
	if !ok {
		return 0, NewErrInvalidType("int32", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsInt32() (int32, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	switch v := av.Value.(type) {
	case int32:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 32)
		return int32(res), err
	}
	return 0, NewErrInvalidType("int32", av.Value)
}

// Int32OrDefault 返回 int32 数据，或者默认值
func (av AnyValue) Int32OrDefault(def int32) int32 {
	val, err := av.Int32()
	if err != nil {
		return def
	}
	return val
}

// Uint32 返回 uint32 数据
func (av AnyValue) Uint32() (uint32, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(uint32)
	if !ok {
		return 0, NewErrInvalidType("uint32", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsUint32() (uint32, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	switch v := av.Value.(type) {
	case uint32:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 32)
		return uint32(res), err
	}
	return 0, NewErrInvalidType("uint32", av.Value)
}

// Uint32OrDefault 返回 uint32 数据，或者默认值
func (av AnyValue) Uint32OrDefault(def uint32) uint32 {
	val, err := av.Uint32()
	if err != nil {
		return def
	}
	return val
}

// Int64 返回 int64 数据
func (av AnyValue) Int64() (int64, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(int64)
	if !ok {
		return 0, NewErrInvalidType("int64", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsInt64() (int64, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	switch v := av.Value.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	}
	return 0, NewErrInvalidType("int64", av.Value)
}

// Int64OrDefault 返回 int64 数据，或者默认值
func (av AnyValue) Int64OrDefault(def int64) int64 {
	val, err := av.Int64()
	if err != nil {
		return def
	}
	return val
}

// Uint64 返回 uint64 数据
func (av AnyValue) Uint64() (uint64, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(uint64)
	if !ok {
		return 0, NewErrInvalidType("uint64", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsUint64() (uint64, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	switch v := av.Value.(type) {
	case uint64:
		return v, nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	}
	return 0, NewErrInvalidType("uint64", av.Value)
}

// Uint64OrDefault 返回 uint64 数据，或者默认值
func (av AnyValue) Uint64OrDefault(def uint64) uint64 {
	val, err := av.Uint64()
	if err != nil {
		return def
	}
	return val
}

// Float32 返回 float32 数据
func (av AnyValue) Float32() (float32, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(float32)
	if !ok {
		return 0, NewErrInvalidType("float32", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsFloat32() (float32, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	switch v := av.Value.(type) {
	case float32:
		return v, nil
	case string:
		res, err := strconv.ParseFloat(v, 32)
		return float32(res), err
	}
	return 0, NewErrInvalidType("float32", av.Value)
}

// Float32OrDefault 返回 float32 数据，或者默认值
func (av AnyValue) Float32OrDefault(def float32) float32 {
	val, err := av.Float32()
	if err != nil {
		return def
	}
	return val
}

// Float64 返回 float64 数据
func (av AnyValue) Float64() (float64, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	val, ok := av.Value.(float64)
	if !ok {
		return 0, NewErrInvalidType("float64", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsFloat64() (float64, error) {
	if av.Error != nil {
		return 0, av.Error
	}
	switch v := av.Value.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	}
	return 0, NewErrInvalidType("float64", av.Value)
}

// Float64OrDefault 返回 float64 数据，或者默认值
func (av AnyValue) Float64OrDefault(def float64) float64 {
	val, err := av.Float64()
	if err != nil {
		return def
	}
	return val
}

// String 返回 string 数据
func (av AnyValue) String() (string, error) {
	if av.Error != nil {
		return "", av.Error
	}
	val, ok := av.Value.(string)
	if !ok {
		return "", NewErrInvalidType("string", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsString() (string, error) {
	if av.Error != nil {
		return "", av.Error
	}

	var val string
	valueOf := reflect.ValueOf(av.Value)
	switch valueOf.Type().Kind() {
	case reflect.String:
		val = valueOf.String()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = strconv.FormatUint(valueOf.Uint(), 10)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = strconv.FormatInt(valueOf.Int(), 10)
	case reflect.Float32:
		val = strconv.FormatFloat(valueOf.Float(), 'f', 10, 32)
	case reflect.Float64:
		val = strconv.FormatFloat(valueOf.Float(), 'f', 10, 64)
	case reflect.Slice:
		if valueOf.Type().Elem().Kind() != reflect.Uint8 {
			return "", NewErrInvalidType("[]byte", av.Value)
		}
		val = string(valueOf.Bytes())
	default:
		return "", errors.New("未兼容类型，暂时无法转换")
	}

	return val, nil
}

// StringOrDefault 返回 string 数据，或者默认值
func (av AnyValue) StringOrDefault(def string) string {
	val, err := av.String()
	if err != nil {
		return def
	}
	return val
}

// Bytes 返回 []byte 数据
func (av AnyValue) Bytes() ([]byte, error) {
	if av.Error != nil {
		return nil, av.Error
	}
	val, ok := av.Value.([]byte)
	if !ok {
		return nil, NewErrInvalidType("[]byte", av.Value)
	}
	return val, nil
}

func (av AnyValue) AsBytes() ([]byte, error) {
	if av.Error != nil {
		return []byte{}, av.Error
	}
	switch v := av.Value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	}

	return []byte{}, NewErrInvalidType("[]byte", av.Value)
}

// BytesOrDefault 返回 []byte 数据，或者默认值
func (av AnyValue) BytesOrDefault(def []byte) []byte {
	val, err := av.Bytes()
	if err != nil {
		return def
	}
	return val
}

// Bool 返回 bool 数据
func (av AnyValue) Bool() (bool, error) {
	if av.Error != nil {
		return false, av.Error
	}
	val, ok := av.Value.(bool)
	if !ok {
		return false, NewErrInvalidType("bool", av.Value)
	}
	return val, nil
}

// BoolOrDefault 返回 bool 数据，或者默认值
func (av AnyValue) BoolOrDefault(def bool) bool {
	val, err := av.Bool()
	if err != nil {
		return def
	}
	return val
}

// NewErrInvalidType 创建一个代表类型转换失败的错误
func NewErrInvalidType(want string, got any) error {
	return fmt.Errorf("value: 类型转换失败，预期类型:%s, 实际值:%#v", want, got)
}
