package utils

import (
	"database/sql/driver"
	"encoding/json"
)

type KV map[string]any

func (kv KV) Scan(value any) error {
	return json.Unmarshal([]byte(value.(string)), &kv)
}

func (kv KV) Value() (driver.Value, error) {
	if len(kv) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(kv)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
