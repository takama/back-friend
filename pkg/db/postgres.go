package db

import (
	"database/sql"
	"fmt"
	"net/url"

	// PostgreSQL driver import
	_ "github.com/lib/pq"

	"github.com/takama/back-friend/pkg/config"
)

// PostgreSQL implements PostgreSQL driver
type PostgreSQL struct {
	pool *sql.DB
}

// NewPostgreSQL creates new PostgreSQL connection
func NewPostgreSQL(cfg *config.Config) (pg *PostgreSQL, name string, err error) {
	name = cfg.DbType
	dsn, err := url.Parse(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DbUsername, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName))
	if err != nil {
		return
	}
	pg = new(PostgreSQL)
	if pg.pool, err = sql.Open("postgres", dsn.String()); err != nil {
		return nil, name, err
	}
	if err := pg.pool.QueryRow("SELECT version()").Scan(&name); err != nil {
		return nil, name, err
	}
	return pg, name, pg.pool.Ping()
}

// Ready returns DB state
func (pg PostgreSQL) Ready() bool {
	if err := pg.pool.Ping(); err == nil {
		return true
	}
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
