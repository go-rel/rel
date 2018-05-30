package changeset

// Options applicable to changeset.
type Options struct {
	message string
	code    int
	name    string
	exact   bool
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

// Code for changeset operation's error.
func Code(code int) Option {
	return func(opts *Options) {
		opts.code = code
	}
}

// Name is used to define index name of constraints.
func Name(name string) Option {
	return func(opts *Options) {
		opts.name = name
	}
}

// Exact is used to define how index name is matched.
func Exact(exact bool) Option {
	return func(opts *Options) {
		opts.exact = exact
	}
}
