package toml

import (
	"github.com/pelletier/go-toml/v2"

	"github.com/apus-run/sea-kit/encoding"
)

type tomlCodec struct{}

func (c tomlCodec) Name() string {
	return "toml"
}

func (c tomlCodec) Marshal(v interface{}) ([]byte, error) {
	return toml.Marshal(v)
}

func (c tomlCodec) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func init() {
	encoding.RegisterCodec(tomlCodec{})
}
