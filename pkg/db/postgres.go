package db

import "github.com/takama/back-friend/pkg/config"

// PostgreSQL implements PostgreSQL driver
type PostgreSQL struct {
}

// NewPostgreSQL creates new PostgreSQL connection
func NewPostgreSQL(cfg *config.Config) (*PostgreSQL, error) {
	pg := new(PostgreSQL)
	return pg, nil
}

// Ready returns DB state
func (pg PostgreSQL) Ready() bool {
	return false
}

// MigrateUp migrates DB schema
func (pg PostgreSQL) MigrateUp() error {
	return nil
}

// MigrateDown remove DB schema and data
func (pg PostgreSQL) MigrateDown() error {
	return nil
}
