package db

import "github.com/takama/back-friend/pkg/config"

// Stub implements Stub driver
type Stub struct {
}

// NewStub creates new Stub structure that implements Driver interface
func NewStub(cfg *config.Config) (stub *Stub, name string, err error) {
	name = cfg.DbType
	stub = new(Stub)
	return
}

// Ready returns DB state
func (stub Stub) Ready() bool {
	return true
}

// MigrateUp migrates DB schema
func (stub Stub) MigrateUp() error {
	return nil
}

// MigrateDown remove DB schema and data
func (stub Stub) MigrateDown() error {
	return nil
}
