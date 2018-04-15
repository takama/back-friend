package db

import (
	"database/sql"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/logger"

	"github.com/takama/backer/datastore"
	"github.com/takama/backer/model"
)

// SQLDriver contains common SQL DB methods
type SQLDriver interface {
	Transaction() (*sql.Tx, error)
	NewPlayer(ID string, tx *sql.Tx) error
	FindPlayer(ID string, tx *sql.Tx) (*model.Player, error)
	SavePlayer(player *model.Player, tx *sql.Tx) error
	NewTournament(ID uint64, tx *sql.Tx) error
	FindTournament(ID uint64, tx *sql.Tx) (*model.Tournament, error)
	SaveTournament(tournament *model.Tournament, tx *sql.Tx) error
}

// Connection implements DB controller interface
type Connection struct {
	*config.Config
	datastore.Store
	datastore.Controller
	SQLDriver
}

// New creates new database connection
func New(cfg *config.Config, log logger.Logger) (conn *Connection, name string, err error) {
	conn = &Connection{Config: cfg}
	pg, json, mock, stub := new(PostgreSQL), new(JSON), new(Mock), new(datastore.Stub)
	switch cfg.DbType {
	case config.PGDriver:
		pg, name, err = NewPostgreSQL(cfg, log)
		conn.SQLDriver, conn.Store = pg, pg
	case config.JSONDriver:
		json, name, err = NewJSON(cfg, log)
		conn.Controller, conn.Store = json, json
	case config.MockDriver:
		conn.Controller, conn.Store, name = mock, mock, config.MockDriver
	default:
		conn.Controller, conn.Store, name = stub, stub, config.StubDriver
	}

	return
}

// Ready returns connection state
func (conn Connection) Ready() bool {
	return conn.Store.Ready()
}

// Transaction returns DB transaction control
func (conn Connection) Transaction() (datastore.Transact, error) {
	if conn.Config.DbType == config.PGDriver {
		return conn.SQLDriver.Transaction()
	}
	return conn.Controller.Transaction()
}

// NewPlayer creates a player with specified ID
func (conn Connection) NewPlayer(ID string, tx datastore.Transact) error {
	if t, ok := tx.(*sql.Tx); ok {
		return conn.SQLDriver.NewPlayer(ID, t)
	}
	return conn.Controller.NewPlayer(ID, tx)
}

// FindPlayer finds existing player by specified ID
func (conn Connection) FindPlayer(ID string, tx datastore.Transact) (*model.Player, error) {
	if t, ok := tx.(*sql.Tx); ok {
		return conn.SQLDriver.FindPlayer(ID, t)
	}
	return conn.Controller.FindPlayer(ID, tx)
}

// SavePlayer saves a Player model
func (conn Connection) SavePlayer(player *model.Player, tx datastore.Transact) error {
	if t, ok := tx.(*sql.Tx); ok {
		return conn.SQLDriver.SavePlayer(player, t)
	}
	return conn.Controller.SavePlayer(player, tx)
}

// NewTournament creates a new tournament with specified ID
func (conn Connection) NewTournament(ID uint64, tx datastore.Transact) error {
	if t, ok := tx.(*sql.Tx); ok {
		return conn.SQLDriver.NewTournament(ID, t)
	}
	return conn.Controller.NewTournament(ID, tx)
}

// FindTournament finds existing tournament by specified ID
func (conn Connection) FindTournament(ID uint64, tx datastore.Transact) (*model.Tournament, error) {
	if t, ok := tx.(*sql.Tx); ok {
		return conn.SQLDriver.FindTournament(ID, t)
	}
	return conn.Controller.FindTournament(ID, tx)
}

// SaveTournament saves a Tournament model
func (conn Connection) SaveTournament(tournament *model.Tournament, tx datastore.Transact) error {
	if t, ok := tx.(*sql.Tx); ok {
		return conn.SQLDriver.SaveTournament(tournament, t)
	}
	return conn.Controller.SaveTournament(tournament, tx)
}

// Reset makes the DB initialization
func (conn Connection) Reset() error {
	return conn.Store.Reset()
}
