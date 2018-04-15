package db

import (
	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/logger"
	"github.com/takama/backer/datastore"
	"github.com/takama/backer/model"
)

// JSON implements JSON files driver
type JSON struct {
}

// NewJSON creates new JSON structure that implements Driver interface
func NewJSON(cfg *config.Config, log logger.Logger) (json *JSON, name string, err error) {
	name = cfg.DbType
	json = new(JSON)
	return
}

// Ready returns DB state
func (json JSON) Ready() bool {
	return true
}

// Reset makes the DB initialization
func (json JSON) Reset() error {
	return nil
}

// MigrateUp migrates DB schema
func (json JSON) MigrateUp() error {
	return nil
}

// MigrateDown remove DB schema and data
func (json JSON) MigrateDown() error {
	return nil
}

// NewPlayer creates a player with specified ID
func (json JSON) NewPlayer(ID string, tx datastore.Transact) error {
	return nil
}

// FindPlayer finds existing player by specified ID
func (json JSON) FindPlayer(ID string, tx datastore.Transact) (*model.Player, error) {
	return nil, nil
}

// SavePlayer saves a Player model
func (json JSON) SavePlayer(player *model.Player, tx datastore.Transact) error {
	return nil
}

// NewTournament creates a new tournament with specified ID
func (json JSON) NewTournament(ID uint64, tx datastore.Transact) error {
	return nil
}

// FindTournament finds existing tournament by specified ID
func (json JSON) FindTournament(ID uint64, tx datastore.Transact) (*model.Tournament, error) {
	return nil, nil
}

// SaveTournament saves a Tournament model
func (json JSON) SaveTournament(tournament *model.Tournament, tx datastore.Transact) error {
	return nil
}

// Transaction returns DB transaction control
func (json JSON) Transaction() (datastore.Transact, error) {
	return new(transactJSON), nil
}

type transactJSON struct{}

func (t transactJSON) Commit() error {
	return nil
}

func (t transactJSON) Rollback() error {
	return nil
}
