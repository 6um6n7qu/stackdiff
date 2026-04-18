package compare

// WithRedact returns a copy of opts with RedactSecrets set.
func WithRedact(opts Options, v bool) Options {
	opts.RedactSecrets = v
	return opts
}

// WithNormalize returns a copy of opts with Normalize set.
func WithNormalize(opts Options, v bool) Options {
	opts.Normalize = v
	return opts
}

// WithKeyPrefix returns a copy of opts with KeyPrefix set.
func WithKeyPrefix(opts Options, prefix string) Options {
	opts.KeyPrefix = prefix
	return opts
}

// WithFilterStatuses returns a copy of opts with FilterStatuses set.
func WithFilterStatuses(opts Options, statuses ...string) Options {
	opts.FilterStatuses = statuses
	return opts
}
