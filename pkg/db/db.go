package db

import "github.com/takama/back-friend/pkg/config"

// Driver contains common DB methods
type Driver interface {
	Ready() bool
	MigrateUp() error
	MigrateDown() error
}

// Connection implements DB controller interface
type Connection struct {
	Driver
}

// New creates new database connection
func New(cfg *config.Config) (conn *Connection, err error) {
	conn = new(Connection)
	switch cfg.DbType {
	case config.PGDriver:
		conn.Driver, err = NewPostgreSQL(cfg)
	case config.JSONDriver:
		conn.Driver, err = NewJSON(cfg)
	default:
		conn.Driver, err = NewStub(cfg)
	}

	return
}
