package db

import (
	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/logger"
)

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
func New(cfg *config.Config, log logger.Logger) (conn *Connection, name string, err error) {
	conn = new(Connection)
	switch cfg.DbType {
	case config.PGDriver:
		conn.Driver, name, err = NewPostgreSQL(cfg, log)
	case config.JSONDriver:
		conn.Driver, name, err = NewJSON(cfg, log)
	case config.MockDriver:
		conn.Driver, name, err = NewMock(cfg, log)
	default:
		conn.Driver, name, err = NewStub(cfg, log)
	}

	return
}

// Ready returns connection state
func (conn Connection) Ready() bool {
	return conn.Driver.Ready()
}
