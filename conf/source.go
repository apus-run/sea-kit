package conf

// KeyValue is conf key value.
type KeyValue struct {
	Key    string
	Value  []byte
	Format string
	Path   string
}

func (k *KeyValue) Read(p []byte) (n int, err error) {
	return copy(p, k.Value), nil
}

type Source interface {
	Load() ([]*KeyValue, error)
}
