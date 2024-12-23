package goconfig

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type serverConfig struct {
	HTTP struct {
		Listen       string        `default:":8080"  json:"listen" yaml:"listen" cli:"http-listen" env:"SERVER_HTTP_LISTEN"`
		ReadTimeout  time.Duration `default:"120s"   json:"read_timeout" yaml:"read_timeout" env:"SERVER_HTTP_READ_TIMEOUT"`
		WriteTimeout time.Duration `default:"120s"   json:"write_timeout" yaml:"write_timeout" env:"SERVER_HTTP_WRITE_TIMEOUT"`
	} `json:"http" yaml:"http"`
	GRPC struct {
		Listen  string        `default:"tcp://:8081" json:"listen" yaml:"listen" cli:"grpc-listen" env:"SERVER_GRPC_LISTEN"`
		Timeout time.Duration `default:"120s"        json:"timeout" yaml:"timeout" env:"SERVER_GRPC_TIMEOUT"`
	} `json:"grpc" yaml:"grpc"`
	Profile struct {
		Mode   string `json:"mode" yaml:"mode" default:"" env:"SERVER_PROFILE_MODE"`
		Listen string `json:"listen" yaml:"listen" default:"" env:"SERVER_PROFILE_LISTEN"`
	} `json:"profile" yaml:"profile"`
}

type testConfig struct {
	configFilepath string
	ServiceName    string `json:"service_name" yaml:"service_name" env:"SERVICE_NAME" default:"disk"`

	LogAddr  string `json:"log_addr" yaml:"log_addr" default:"" env:"LOG_ADDR"`
	LogLevel string `json:"log_level" yaml:"log_level" default:"debug" env:"LOG_LEVEL" cli:"log-level" short-cli:"l"`

	Server serverConfig `json:"server" yaml:"server"`
}

func (cfg *testConfig) ConfigFilepath() string {
	return cfg.configFilepath
}

func TestConfigLoadEnv(t *testing.T) {
	var conf testConfig

	os.Args = []string{"test", "--http-listen=addr:test", "-l", "error"}
	os.Setenv("SERVICE_NAME", "test-servername")
	os.Setenv("LOG_ADDR", "test-logger-addr")
	os.Setenv("LOG_LEVEL", "error-loglevel")

	assert.NoError(t, Load(&conf, WithDefaults(), WithEnv(), WithArgs()))
	assert.Equal(t, "test-servername", conf.ServiceName)
	assert.Equal(t, "test-logger-addr", conf.LogAddr)
	assert.Equal(t, "error", conf.LogLevel)
	assert.Equal(t, "addr:test", conf.Server.HTTP.Listen)
}

func TestConfigLoadWithEnvOptsOnly(t *testing.T) {
	var conf testConfig

	os.Setenv("SERVICE_NAME", "test-servername")
	os.Setenv("LOG_ADDR", "test-logger-addr")
	os.Setenv("LOG_LEVEL", "error-loglevel")

	assert.NoError(t, Load(&conf, WithEnv(), WithCustomArgs("--http-listen", "addr:test")))
	assert.Equal(t, "test-servername", conf.ServiceName)
	assert.Equal(t, "test-logger-addr", conf.LogAddr)
	assert.Equal(t, "error-loglevel", conf.LogLevel)
}

func TestConfigLoadEnvError(t *testing.T) {
	var conf testConfig

	os.Args = []string{}
	os.Setenv("SERVER_HTTP_READ_TIMEOUT", "error")

	assert.Error(t, Load(&conf, WithEnv(), WithDefaults()))
}

func TestConfigLoadCliError(t *testing.T) {
	var conf testConfig
	os.Args = []string{"test", "-v"}
	assert.Error(t, Load(&conf))
}

func TestConfigLoadFile(t *testing.T) {
	configs := []string{
		"test-assets/config.hcl",
		"test-assets/config.json",
		"test-assets/config.yml",
	}

	// Reset CLI and ENV parameters
	os.Args = []string{}
	os.Setenv("SERVICE_NAME", "")
	os.Setenv("LOG_ADDR", "")
	os.Setenv("LOG_LEVEL", "")
	os.Setenv("SERVER_HTTP_LISTEN", "")
	os.Setenv("SERVER_GRPC_LISTEN", "")
	os.Setenv("SERVER_PROFILE_MODE", "")
	os.Setenv("SERVER_HTTP_READ_TIMEOUT", "")

	for _, configFile := range configs {
		var conf testConfig
		conf.configFilepath = configFile
		assert.NoError(t, Load(&conf, WithDefaults(), WithArgs(), WithEnv(), WithFile("")))

		assert.Equal(t, "test-servername", conf.ServiceName)
		assert.Equal(t, "logstash", conf.LogAddr)
		assert.Equal(t, "error", conf.LogLevel)
		assert.Equal(t, "test:port", conf.Server.HTTP.Listen)
		assert.Equal(t, "test:port", conf.Server.GRPC.Listen)
		assert.Equal(t, "net", conf.Server.Profile.Mode)
	}
}

func TestConfigLoadFileUnsupported(t *testing.T) {
	// Reset CLI and ENV parameters
	os.Args = []string{}
	os.Setenv("SERVER_HTTP_READ_TIMEOUT", "")

	var conf testConfig
	assert.Error(t, Load(&conf, WithFile("test-assets/config.unsupported")))
}

func TestConfigLoadFileOpenError(t *testing.T) {
	// Reset CLI and ENV parameters
	os.Args = []string{}
	os.Setenv("SERVER_HTTP_READ_TIMEOUT", "")

	var conf testConfig
	conf.configFilepath = "undefined"
	assert.Error(t, Load(&conf, WithFile("")))
}
