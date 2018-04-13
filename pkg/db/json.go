package db

import "github.com/takama/back-friend/pkg/config"

// JSON implements JSON files driver
type JSON struct {
}

// NewJSON creates new JSON structure that implements Driver interface
func NewJSON(cfg *config.Config) (json *JSON, name string, err error) {
	name = cfg.DbType
	json = new(JSON)
	return
}

// Ready returns DB state
func (json JSON) Ready() bool {
	return true
}

// MigrateUp migrates DB schema
func (json JSON) MigrateUp() error {
	return nil
}

// MigrateDown remove DB schema and data
func (json JSON) MigrateDown() error {
	return nil
}
