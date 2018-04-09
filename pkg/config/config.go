package config

import "github.com/kelseyhightower/envconfig"

// Default values: path to config file, host, port, etc
const (
	// ServiceName contains a service name prefix which used in ENV variables
	ServiceName = "backfriend"
)

// Config - Service configuration
type Config struct {
}

// New - returns new config record initialized with default values
func New() *Config {
	return &Config{}
}

// LoadFromEnv load configuration parameters from environment
func (config *Config) LoadFromEnv() error {

	// Load all rest environment
	return envconfig.Process(ServiceName, config)
}
