# GoConfig

[![Build Status](https://github.com/demdxx/goconfig/workflows/run%20tests/badge.svg)](https://github.com/demdxx/goconfig/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/goconfig)](https://goreportcard.com/report/github.com/demdxx/goconfig)
[![GoDoc](https://godoc.org/github.com/demdxx/goconfig?status.svg)](https://godoc.org/github.com/demdxx/goconfig)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/goconfig/badge.svg)](https://coveralls.io/github/demdxx/goconfig)

Goconfig is a Go (Golang) configuration initialization module that provides simple and efficient functionality for loading configurations based on struct definitions. It supports multiple configuration sources such as defaults, environment variables, command-line arguments, and configuration files in various formats (JSON, YAML, HCL).

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Defining Configuration Structures](#defining-configuration-structures)
  - [Loading Configuration](#loading-configuration)
    - [Loading with Default Options](#loading-with-default-options)
    - [Loading with Specific Options](#loading-with-specific-options)
- [Example](#example)
  - [config.go](#configgo)
  - [main.go](#maingo)
- [Dependencies](#dependencies)
- [Contributing](#contributing)
- [License](#license)
- [TODO](#todo)

## Features

- Struct-based Configuration: Define your configuration using Go structs with tags for JSON, YAML, CLI, and environment variables.
- Multiple Sources: Load configuration from defaults, environment variables, command-line arguments, and configuration files.
- Flexible File Support: Supports JSON, YAML, and HCL configuration file formats.
- Extensible: Easily extendable to support additional configuration sources or formats.

## Installation

```sh
go get github.com/demdxx/goconfig
```

## Usage

### Defining Configuration Structures

Define your configuration using Go structs. Use struct tags to specify default values and mappings for different configuration sources.

```go
package config

import "time"

type ServerConfig struct {
 HTTP struct {
  Listen       string        `default:":8080"  json:"listen" yaml:"listen" cli:"http-listen" env:"SERVER_HTTP_LISTEN"`
  ReadTimeout  time.Duration `default:"120s"   json:"read_timeout" yaml:"read_timeout" env:"SERVER_HTTP_READ_TIMEOUT"`
  WriteTimeout time.Duration `default:"120s"   json:"write_timeout" yaml:"write_timeout" env:"SERVER_HTTP_WRITE_TIMEOUT"`
 }
 GRPC struct {
  Listen  string        `default:"tcp://:8081" json:"listen" yaml:"listen" cli:"grpc-listen" env:"SERVER_GRPC_LISTEN"`
  Timeout time.Duration `default:"120s"        json:"timeout" yaml:"timeout" env:"SERVER_GRPC_TIMEOUT"`
 }
 Profile struct {
  Mode   string `json:"mode" yaml:"mode" default:"" env:"SERVER_PROFILE_MODE"`
  Listen string `json:"listen" yaml:"listen" default:"" env:"SERVER_PROFILE_LISTEN"`
 }
}

type ConfigType struct {
 ServiceName     string       `json:"service_name" yaml:"service_name" env:"SERVICE_NAME" default:"disk"`
 DatacenterName  string       `json:"datacenter_name" yaml:"datacenter_name" env:"DC_NAME" default:"??"`
 Hostname        string       `json:"hostname" yaml:"hostname" env:"HOSTNAME" default:""`
 Hostcode        string       `json:"hostcode" yaml:"hostcode" env:"HOSTCODE" default:""`

 LogAddr         string       `default:"" env:"LOG_ADDR"`
 LogLevel        string       `default:"debug" env:"LOG_LEVEL"`

 Server          ServerConfig `json:"server" yaml:"server"`
}

var Config ConfigType
```

### Loading Configuration

Use the Load function to initialize your configuration. You can specify various options such as loading defaults, environment variables, command-line arguments, and configuration files.

#### Loading with Default Options

By default, Load will attempt to load configuration from defaults, environment variables, command-line arguments, and configuration files.

```go
package main

import (
 "log"

 configLoader "github.com/demdxx/goconfig"
 "your_project/config"
)

func init() {
 if err := configLoader.Load(&config.Config); err != nil {
  log.Fatalf("Failed to load configuration: %v", err)
 }
}

func main() {
 // Your application code
}
```

#### Loading with Specific Options

You can customize the loading process by specifying options such as WithDefaults, WithEnv, WithArgs, and WithFile.

```go
package main

import (
 "log"

 "github.com/demdxx/goconfig"
 "your_project/config"
)

func init() {
 // Example: Load configuration with defaults and environment variables only
 options := []goconfig.Option{
  goconfig.WithDefaults(),
  goconfig.WithEnv(),
 }

 if err := goconfig.Load(&config.Config, options...); err != nil {
  log.Fatalf("Failed to load configuration: %v", err)
 }
}

func main() {
 // Your application code
}
```

Available Options:

- `WithDefaults()`: Sets default values for the configuration.
- `WithEnv()`: Parses environment variables.
- `WithArgs(...string)`: Parses command-line arguments.
- `WithFile(path string)`: Loads configuration from a specified file.

## Example

### config.go

```go
package config

import "time"

type ServerConfig struct {
 HTTP struct {
  Listen       string        `default:":8080"  json:"listen" yaml:"listen" cli:"http-listen" env:"SERVER_HTTP_LISTEN"`
  ReadTimeout  time.Duration `default:"120s"   json:"read_timeout" yaml:"read_timeout" env:"SERVER_HTTP_READ_TIMEOUT"`
  WriteTimeout time.Duration `default:"120s"   json:"write_timeout" yaml:"write_timeout" env:"SERVER_HTTP_WRITE_TIMEOUT"`
 }
 GRPC struct {
  Listen  string        `default:"tcp://:8081" json:"listen" yaml:"listen" cli:"grpc-listen" env:"SERVER_GRPC_LISTEN"`
  Timeout time.Duration `default:"120s"        json:"timeout" yaml:"timeout" env:"SERVER_GRPC_TIMEOUT"`
 }
 Profile struct {
  Mode   string `json:"mode" yaml:"mode" default:"" env:"SERVER_PROFILE_MODE"`
  Listen string `json:"listen" yaml:"listen" default:"" env:"SERVER_PROFILE_LISTEN"`
 }
}

type ConfigType struct {
 ServiceName     string       `json:"service_name" yaml:"service_name" env:"SERVICE_NAME" default:"disk"`
 DatacenterName  string       `json:"datacenter_name" yaml:"datacenter_name" env:"DC_NAME" default:"??"`
 Hostname        string       `json:"hostname" yaml:"hostname" env:"HOSTNAME" default:""`
 Hostcode        string       `json:"hostcode" yaml:"hostcode" env:"HOSTCODE" default:""`

 LogAddr         string       `default:"" env:"LOG_ADDR"`
 LogLevel        string       `default:"debug" env:"LOG_LEVEL"`

 Server          ServerConfig `json:"server" yaml:"server"`
}

var Config ConfigType
```

### main.go

```go
package main

import (
 "log"

 "github.com/demdxx/goconfig"
 "your_project/config"
)

func init() {
 // Load configuration with defaults, environment variables, and a specific config file
 options := []goconfig.Option{
  goconfig.WithDefaults(),
  goconfig.WithEnv(),
  goconfig.WithFile("config.yaml"),
 }

 if err := goconfig.Load(&config.Config, options...); err != nil {
  log.Fatalf("Failed to load configuration: %v", err)
 }
}

func main() {
 // Your application code
}
```

## Dependencies

- github.com/caarlos0/env
- github.com/hashicorp/hcl
- github.com/mcuadros/go-defaults

## Contributing

Contributions are welcome! Please open issues and submit pull requests for any features, bug fixes, or improvements.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## TODO

- [ ] Add support for environment variable prefixes
