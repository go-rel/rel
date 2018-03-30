package changeset

type Options struct {
	Message    string
	AllowError bool
}

type Option func(*Options)

func (options *Options) Apply(opts []Option) {
	for _, opt := range opts {
		opt(options)
	}
}

func Message(message string) Option {
	return func(opts *Options) {
		opts.Message = message
	}
}

func AllowError(allow bool) Option {
	return func(opts *Options) {
		opts.AllowError = allow
	}
}
