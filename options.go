package goconfig

import "os"

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
func WithArgs() Option {
	return func(o *options) {
		args := os.Args
		if len(args) > 0 {
			args = args[1:]
		}
		o.parsers = append(o.parsers, argsParser(args...))
	}
}

// WithCustomArgs parse custom arguments
func WithCustomArgs(args ...string) Option {
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
