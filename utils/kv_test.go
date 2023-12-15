package utils

import (
	"strconv"
	"testing"
)

func TestKV(t *testing.T) {
	type Foo struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Extra KV     `json:"extra"`
	}
	kv := &Foo{

		Extra: map[string]any{
			"from_id":  strconv.Itoa(1),
			"from_oid": strconv.FormatInt(1, 10),
		},
	}

	kv.Extra.Scan(`{"id": 1, "name": "moocss"}`)
	v, err := kv.Extra.Value()
	if err != nil {
		t.Errorf("格式错误: %v", err)
	}
	t.Logf("输出: %v", v)
}
