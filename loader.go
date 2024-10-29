package goconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	env "github.com/caarlos0/env/v6"
	"github.com/gravitational/configure"
	"github.com/hashicorp/hcl"
	defaults "github.com/mcuadros/go-defaults"
)

type configFilepath interface {
	ConfigFilepath() string
}

type config any

type options struct {
	defaults bool
	args     bool
	file     bool
	env      bool
}

func WithDefaults() func(*options) {
	return func(o *options) {
		o.defaults = true
	}

}

func WithArgs() func(*options) {
	return func(o *options) {
		o.args = true
	}
}

func WithFile() func(*options) {
	return func(o *options) {
		o.file = true
	}
}

func WithEnv() func(*options) {
	return func(o *options) {
		o.env = true
	}
}

// Load data from file
func Load(cfg config, opts ...func(*options)) (err error) {
	o := &options{}
	if opts == nil {
		o.env = true
		o.args = true
		o.file = true
		o.defaults = true
	} else {
		for _, opt := range opts {
			opt(o)
		}
	}

	// Set defaults for config
	if o.defaults {
		defaults.SetDefaults(cfg)
	}

	// parse command line arguments
	if o.args {
		if len(os.Args) > 1 {
			if err = configure.ParseCommandLine(cfg, os.Args[1:]); err != nil {
				return err
			}
		}
	}

	// parse config from file
	if o.file {
		if configFile, _ := cfg.(configFilepath); configFile != nil {
			if filepath := configFile.ConfigFilepath(); len(filepath) > 0 {
				if err = loadFile(cfg, filepath); err != nil {
					return err
				}
			}
		}
	}

	// parse environment variables
	if o.env {
		if err = env.Parse(cfg); err != nil {
			return err
		}
	}

	return err
}

// loadFile config from file path
func loadFile(cfg config, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".yml", ".yaml":
		return configure.ParseYAML(data, cfg)
	case ".json":
		return json.Unmarshal(data, cfg)
	case ".hcl":
		var root any
		// For some specific HCL module not always work as expected
		// so this is a hack to fix it
		if err = hcl.Unmarshal(data, &root); err != nil {
			return err
		}
		if data, err = json.Marshal(root); err != nil {
			return err
		}
		// Skip the error because of HCL converts structures into arrays of structs
		_ = json.Unmarshal(data, cfg)
		return hcl.Unmarshal(data, cfg)
	}
	return fmt.Errorf("unsupported config ext: %s", ext)
}
