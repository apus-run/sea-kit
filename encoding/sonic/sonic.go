package sonic

import (
	"github.com/bytedance/sonic"

	"github.com/apus-run/sea-kit/encoding"
)

// Name is the name registered for the msgpack compressor.
const Name = "sonic"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with msgpack.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
