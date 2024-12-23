package goconfig

type (
	extParser func(cfg config) error
	Option    func(*options)
)

type options struct {
	parsers []extParser
}

// WithDefaults set defaults for config
func WithDefaults() Option {
	return func(o *options) {
		o.parsers = append(o.parsers, defaultsParser)
	}
}

// WithArgs parse command line arguments
func WithArgs(args ...string) Option {
	return func(o *options) {
		o.parsers = append(o.parsers, argsParser(args...))
	}
}

// WithFile parse config from file
func WithFile(path string) Option {
	return func(o *options) {
		o.parsers = append(o.parsers, fileParser(path))
	}
}

// WithEnv parse environment variables
func WithEnv() Option {
	return func(o *options) {
		o.parsers = append(o.parsers, envParser)
	}
}
