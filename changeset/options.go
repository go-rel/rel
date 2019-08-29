package changeset

// Options applicable to changeset.
type Options struct {
	message     string
	key         string
	exact       bool
	changeOnly  bool
	required    bool
	sourceField string
	emptyValues []interface{}
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

// Key is used to define index key of constraints.
func Key(key string) Option {
	return func(opts *Options) {
		opts.key = key
	}
}

// Exact is used to define how index key is matched.
func Exact(exact bool) Option {
	return func(opts *Options) {
		opts.exact = exact
	}
}

// ChangeOnly is used to define if validate is only check change
func ChangeOnly(changeOnly bool) Option {
	return func(opts *Options) {
		opts.changeOnly = changeOnly
	}
}

// Required is used to define whether an assoc needs to be required or not.
func Required(required bool) Option {
	return func(opts *Options) {
		opts.required = required
	}
}

// SourceField to define used field name in params.
func SourceField(field string) Option {
	return func(opts *Options) {
		opts.sourceField = field
	}
}

// EmptyValues defines list of empty values when casting.
// default to [""]
func EmptyValues(values ...interface{}) Option {
	return func(opts *Options) {
		opts.emptyValues = values
	}
}
