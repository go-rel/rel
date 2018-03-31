package changeset

type Options struct {
	message    string
	allowError bool
}

type Option func(*Options)

func (options *Options) Apply(opts []Option) {
	for _, opt := range opts {
		opt(options)
	}
}

func Message(message string) Option {
	return func(opts *Options) {
		opts.message = message
	}
}

func AllowError(allow bool) Option {
	return func(opts *Options) {
		opts.allowError = allow
	}
}
