package changeset

type Options struct {
	Message string
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
