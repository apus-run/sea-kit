package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ReverseMap just reverse the map from [key, value] to [value, key]
func ReverseMap[K comparable, V comparable](m map[K]V) map[V]K {
	n := make(map[V]K, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

// EnumMarshalYaml is a help function to marshal the enum to yaml
func EnumMarshalYaml[T comparable](m map[T]string, v T, typename string) (interface{}, error) {
	if val, ok := m[v]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("%v is not a valid %s", v, typename)
}

// EnumMarshalJSON is a help function to marshal the enum to JSON
func EnumMarshalJSON[T comparable](m map[T]string, v T, typename string) ([]byte, error) {
	if val, ok := m[v]; ok {
		return []byte(fmt.Sprintf(`"%s"`, val)), nil
	}
	return nil, fmt.Errorf("%v is not a valid %s", v, typename)
}

// EnumUnmarshalYaml is a help function to unmarshal the enum from yaml
func EnumUnmarshalYaml[T comparable](unmarshal func(interface{}) error, m map[string]T, v *T, init T, typename string) error {
	var str string
	*v = init
	if err := unmarshal(&str); err != nil {
		return err
	}
	if val, ok := m[strings.ToLower(str)]; ok {
		*v = val
		return nil
	}
	return fmt.Errorf("%v is not a valid %s", str, typename)
}

// EnumUnmarshalJSON is a help function to unmarshal the enum from JSON
func EnumUnmarshalJSON[T comparable](b []byte, m map[string]T, v *T, init T, typename string) error {
	var str string
	*v = init
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	if val, ok := m[strings.ToLower(str)]; ok {
		*v = val
		return nil
	}
	return fmt.Errorf("%v is not a valid %s", str, typename)
}
