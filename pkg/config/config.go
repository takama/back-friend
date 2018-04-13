package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/takama/back-friend/pkg/helper"
	"github.com/takama/back-friend/pkg/logger"
)

// Default values: host, port, etc
const (
	// ServiceName contains a service name prefix which used in ENV variables
	ServiceName = "BackFriend"

	defaultHost       = "0.0.0.0"
	defaultPort       = 7117
	defaultLogLevel   = logger.LevelInfo
	defaultDBLocation = "/var/lib/datastore"

	// StubDriver defines driver name for testing
	StubDriver = "stub"

	// JSONDriver defines driver name for JSON files DB
	JSONDriver = "json"

	// PGDriver defines driver name for PostgreSQL DB
	PGDriver = "postgresql"
)

// Config - Service configuration
type Config struct {
	// Local service host
	LocalHost string `split_words:"true"`
	// Local service port
	LocalPort int `split_words:"true"`
	// Logging level in logger.Level notation
	LogLevel logger.Level `split_words:"true"`
	// Database type
	DbType string `split_words:"true"`
	// Database files location
	DbLocation string `split_words:"true"`
	// Database host
	DbHost string `split_words:"true"`
	// Database port
	DbPort int `split_words:"true"`
	// Database name
	DbName string `split_words:"true"`
	// Database username
	DbUsername string `split_words:"true"`
	// Database password
	DbPassword string `split_words:"true"`
}

// New - returns new config record initialized with default values
func New() *Config {
	return &Config{
		LocalHost:  defaultHost,
		LocalPort:  defaultPort,
		LogLevel:   defaultLogLevel,
		DbType:     StubDriver,
		DbLocation: defaultDBLocation,
	}
}

// LoadFromEnv load configuration parameters from environment
func (config *Config) LoadFromEnv() error {

	return envconfig.Process(helper.ToSnake(ServiceName), config)
}
