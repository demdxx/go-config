package goconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	env "github.com/caarlos0/env/v6"
	"github.com/demdxx/gocast/v2"
	"github.com/hashicorp/hcl"
	defaults "github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

type configFilepath interface {
	ConfigFilepath() string
}

type config any

// Load data from file
func Load(cfg config, opts ...Option) error {
	var confOpt options

	// Apply options
	for _, opt := range opts {
		opt(&confOpt)
	}

	// Apply defaults if needed
	if len(confOpt.parsers) == 0 {
		confOpt.parsers = []extParser{defaultsParser, argsParser(), fileParser(""), envParser}
	}

	// Apply parsers
	for _, parser := range confOpt.parsers {
		if err := parser(cfg); err != nil {
			return err
		}
	}

	return nil
}

// defaultsParser set defaults for config
func defaultsParser(cfg config) error {
	defaults.SetDefaults(cfg)
	return nil
}

// argsParser parse command line arguments
func argsParser(args ...string) extParser {
	return func(cfg config) error {
		if len(args) == 0 && len(os.Args) > 1 {
			args = os.Args[1:]
		}
		flags, err := parseCommandFlags(args)
		if err != nil {
			return err
		}
		if len(flags) == 0 {
			return nil
		}
		return gocast.StructWalk(context.Background(), cfg, func(ctx context.Context, _ gocast.StructWalkObject, field gocast.StructWalkField, _ []string) error {
			if cliField := field.Tag("cli"); cliField != "" {
				if val, ok := flags[cliField]; ok {
					return field.SetValue(ctx, val)
				}
			}
			if cliField := field.Tag("short-cli"); cliField != "" {
				if val, ok := flags[cliField]; ok {
					return field.SetValue(ctx, val)
				}
			}
			return nil
		})
	}
}

// fileParser parse config from file
func fileParser(path string) extParser {
	return func(cfg config) error {
		if path != "" {
			return loadFile(cfg, path)
		} else if configFile, _ := cfg.(configFilepath); configFile != nil {
			if filepath := configFile.ConfigFilepath(); len(filepath) > 0 {
				return loadFile(cfg, filepath)
			}
		}
		return nil
	}
}

// envParser parse environment variables
func envParser(cfg config) error {
	return env.Parse(cfg)
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
		// return configure.ParseYAML(data, cfg)
		return yaml.Unmarshal(data, cfg)
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
