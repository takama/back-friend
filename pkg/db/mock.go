package db

import (
	"github.com/takama/backer/datastore"
	"github.com/takama/backer/model"
)

// Mock implements Mock driver
type Mock struct {
	OnReady          func() bool
	OnReset          func() error
	OnMigrateUp      func() error
	OnMigrateDown    func() error
	OnNewPlayer      func(ID string, tx datastore.Transact) error
	OnFindPlayer     func(ID string, tx datastore.Transact) (*model.Player, error)
	OnSavePlayer     func(player *model.Player, tx datastore.Transact) error
	OnNewTournament  func(ID uint64, tx datastore.Transact) error
	OnFindTournament func(ID uint64, tx datastore.Transact) (*model.Tournament, error)
	OnSaveTournament func(tournament *model.Tournament, tx datastore.Transact) error
}

// Ready returns datastore state
func (mock Mock) Ready() bool {
	return mock.OnReady()
}

// Reset makes the DB initialization
func (mock Mock) Reset() error {
	return mock.OnReset()
}

// NewPlayer creates a player with specified ID
func (mock Mock) NewPlayer(ID string, tx datastore.Transact) error {
	return mock.OnNewPlayer(ID, tx)
}

// FindPlayer finds existing player by specified ID
func (mock Mock) FindPlayer(ID string, tx datastore.Transact) (*model.Player, error) {
	return mock.OnFindPlayer(ID, tx)
}

// SavePlayer saves a Player model
func (mock Mock) SavePlayer(player *model.Player, tx datastore.Transact) error {
	return mock.OnSavePlayer(player, tx)
}

// NewTournament creates a new tournament with specified ID
func (mock Mock) NewTournament(ID uint64, tx datastore.Transact) error {
	return mock.OnNewTournament(ID, tx)
}

// FindTournament finds existing tournament by specified ID
func (mock Mock) FindTournament(ID uint64, tx datastore.Transact) (*model.Tournament, error) {
	return mock.OnFindTournament(ID, tx)
}

// SaveTournament saves a Tournament model
func (mock Mock) SaveTournament(tournament *model.Tournament, tx datastore.Transact) error {
	return mock.OnSaveTournament(tournament, tx)
}

// MigrateUp migrates DB schema
func (mock Mock) MigrateUp() error {
	return mock.OnMigrateUp()
}

// MigrateDown remove DB schema and data
func (mock Mock) MigrateDown() error {
	return mock.OnMigrateDown()
}

// Transaction returns DB transaction control
func (mock Mock) Transaction() (datastore.Transact, error) {
	return new(transactMock), nil
}

type transactMock struct{}

func (t transactMock) Commit() error {
	return nil
}

func (t transactMock) Rollback() error {
	return nil
}
