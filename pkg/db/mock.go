package db

import (
	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/logger"
)

// Mock implements Mock driver
type Mock struct {
	OnReady       func() bool
	OnMigrateUp   func() error
	OnMigrateDown func() error
}

// NewMock creates new Mock structure that implements Driver interface
func NewMock(cfg *config.Config, log logger.Logger) (mock *Mock, name string, err error) {
	name = cfg.DbType
	mock = new(Mock)
	return
}

// Ready returns DB state
func (mock Mock) Ready() bool {
	return mock.OnReady()
}

// MigrateUp migrates DB schema
func (mock Mock) MigrateUp() error {
	return mock.OnMigrateUp()
}

// MigrateDown remove DB schema and data
func (mock Mock) MigrateDown() error {
	return mock.OnMigrateDown()
}
