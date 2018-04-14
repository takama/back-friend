package db

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/takama/back-friend/pkg/config"

	"github.com/lib/pq"
)

const (
	// ErrorCodeDuplicateDatabase is code when create database but it already exists
	ErrorCodeDuplicateDatabase = "42P04"
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
	if _, err := pg.pool.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DbName)); err != nil {
		if dbErr, ok := err.(*pq.Error); ok && dbErr.Code == ErrorCodeDuplicateDatabase {
			// Database exists, there is no need to use migration
			return pg, name, pg.pool.Ping()
		}
		return nil, name, err
	}
	// Database has been created, it should migrate up
	if err = pg.MigrateUp(); err != nil {
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
	_, err := pg.pool.Exec(migrateUpData)
	return err
}

// MigrateDown remove DB schema and data
func (pg PostgreSQL) MigrateDown() error {
	_, err := pg.pool.Exec(migrateDownData)
	return err
}

const migrateUpData = `
CREATE TABLE IF NOT EXISTS players (
	id 				VARCHAR(128) PRIMARY KEY,
	balance			NUMERIC(12,2) DEFAULT 0 CHECK (balance >= 0) NOT NULL
);
CREATE TABLE IF NOT EXISTS tournaments (
	id 				INT PRIMARY KEY,
	is_finished 	BOOLEAN DEFAULT FALSE NOT NULL,
	deposit			NUMERIC(12,2) DEFAULT 0 CHECK (deposit >= 0) NOT NULL
);
CREATE TABLE IF NOT EXISTS bidders (
	player_id		VARCHAR(128) REFERENCES players (id),
	tournament_id	INT REFERENCES tournaments (id),
	winner			BOOLEAN DEFAULT FALSE NOT NULL,
	prize			NUMERIC(12,2) DEFAULT 0 CHECK (prize >= 0) NOT NULL,
	backers			JSONB,
	PRIMARY KEY (player_id, tournament_id)
);
`

const migrateDownData = `
DROP TABLE IF EXISTS bidders;
DROP TABLE IF EXISTS tournaments;
DROP TABLE IF EXISTS players;
`
