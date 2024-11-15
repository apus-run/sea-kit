package filesystems

import (
	"io/fs"
)

type Config struct {
	fsmap FileSystems
}

func NewConfig(fsmap FileSystems) *Config {
	return &Config{fsmap: fsmap}
}

func (c *Config) FileSystems() FileSystems {
	return c.fsmap
}

func (c *Config) SetFileSystems(fsmap FileSystems) {
	c.fsmap = fsmap
}

func (c *Config) Register(k string, v fs.FS) {
	c.fsmap.Register(k, v)
}

func (c *Config) Unregister(k string) {
	c.fsmap.Unregister(k)
}

func (c *Config) Get(k string) (v fs.FS, ok bool) {
	return c.fsmap.Get(k)
}

func (c *Config) Default() fs.FS {
	return c.fsmap.Default()
}
