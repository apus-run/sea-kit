package conf

// Option is config option.
type Option func(*Options)

type Options struct {
	sources []Source
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		sources: []Source{},
	}
}

// WithSource with config source.
func WithSource(s ...Source) Option {
	return func(o *Options) {
		o.sources = s
	}
}
