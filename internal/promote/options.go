package promote

// Option is a functional option for configuring a Promoter.
type Option func(*Config)

// Config controls promotion behaviour.
type Config struct {
	// DryRun prevents any mutations; ops are computed but not applied.
	DryRun bool
	// IgnoreKeys is a set of exact keys that should be excluded from promotion.
	IgnoreKeys map[string]struct{}
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		IgnoreKeys: make(map[string]struct{}),
	}
}

// WithDryRun sets the DryRun flag.
func WithDryRun(v bool) Option {
	return func(c *Config) {
		c.DryRun = v
	}
}

// WithIgnoreKeys registers keys that must not be promoted.
func WithIgnoreKeys(keys ...string) Option {
	return func(c *Config) {
		for _, k := range keys {
			c.IgnoreKeys[k] = struct{}{}
		}
	}
}

// NewConfig builds a Config from the supplied options.
func NewConfig(opts ...Option) Config {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
