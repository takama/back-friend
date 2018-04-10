package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/takama/back-friend/pkg/helper"
)

// Default values: host, port, etc
const (
	// ServiceName contains a service name prefix which used in ENV variables
	ServiceName = "BackFriend"

	defaultHost = "0.0.0.0"
	defaultPort = 8080
)

// Config - Service configuration
type Config struct {
	// Local service host
	LocalHost string `split_words:"true"`
	// Local service port
	LocalPort int `split_words:"true"`
}

// New - returns new config record initialized with default values
func New() *Config {
	return &Config{
		LocalHost: defaultHost,
		LocalPort: defaultPort,
	}
}

// LoadFromEnv load configuration parameters from environment
func (config *Config) LoadFromEnv() error {

	return envconfig.Process(helper.ToSnake(ServiceName), config)
}
