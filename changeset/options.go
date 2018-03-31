package changeset

// Options applicable to changeset.
type Options struct {
	message string
}

// Option for changeset operation.
type Option func(*Options)

func (options *Options) apply(opts []Option) {
	for _, opt := range opts {
		opt(options)
	}
}

// Message for changeset operation's error.
func Message(message string) Option {
	return func(opts *Options) {
		opts.message = message
	}
}
