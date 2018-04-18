package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/logger"

	"github.com/lib/pq"
	"github.com/takama/backer/model"
)

// PostgreSQL implements PostgreSQL driver
type PostgreSQL struct {
	logger.Logger
	pool *sql.DB
}

// NewPostgreSQL creates new PostgreSQL connection
func NewPostgreSQL(cfg *config.Config, log logger.Logger) (pg *PostgreSQL, name string, err error) {
	name = cfg.DbType
	dsn, err := url.Parse(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DbUsername, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName))
	if err != nil {
		return
	}
	pg = new(PostgreSQL)
	pg.Logger = log
	if pg.pool, err = sql.Open("postgres", dsn.String()); err != nil {
		return nil, name, err
	}
	if err := pg.pool.QueryRow("SELECT version()").Scan(&name); err != nil {
		return nil, name, err
	}
	if _, err := pg.pool.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DbName)); err != nil {
		if dbErr, ok := err.(*pq.Error); ok && dbErr.Code == errorCodeDuplicateDatabase {
			// Database exists, there is no need to use migration
			return pg, name, pg.pool.Ping()
		}
		return nil, name, err
	}
	log.Infof("Database `%s` has been created. Running migration...", cfg.DbName)
	if err = pg.MigrateUp(); err != nil {
		return nil, name, err
	}
	log.Info("Migration has been done")
	return pg, name, pg.pool.Ping()
}

// Ready returns DB state
func (pg PostgreSQL) Ready() bool {
	if err := pg.pool.Ping(); err == nil {
		return true
	}
	return false
}

// Reset makes the DB initialization
func (pg PostgreSQL) Reset() error {
	var err error
	tx, err := pg.Transaction()
	if err != nil {
		return err
	}
	if err = pg.MigrateDown(); err != nil {
		tx.Rollback()
		return err
	}
	if err = pg.MigrateUp(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// NewPlayer creates a player with specified ID
func (pg PostgreSQL) NewPlayer(ID string, tx *sql.Tx) error {
	stmt, err := tx.Prepare(queryInsertPlayer)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}
	pg.Logger.Debugf("Created a new player %s", ID)
	return err
}

// FindPlayer finds existing player by specified ID
func (pg PostgreSQL) FindPlayer(ID string, tx *sql.Tx) (*model.Player, error) {
	row := tx.QueryRow(queryGetPlayerBalance, ID)
	player := &model.Player{ID: ID}
	err := row.Scan(&player.Balance)
	return player, err
}

// SavePlayer saves a Player model
func (pg PostgreSQL) SavePlayer(player *model.Player, tx *sql.Tx) error {
	stmt, err := tx.Prepare(queryUpdatePlayer)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(player.Balance, player.ID)
	if err != nil {
		return err
	}
	pg.Logger.Debugf("Saved the player %s", player.ID)
	return err
}

// NewTournament creates a new tournament with specified ID
func (pg PostgreSQL) NewTournament(ID uint64, tx *sql.Tx) error {
	stmt, err := tx.Prepare(queryInsertTournament)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}
	pg.Logger.Debugf("Created a new tournament %d", ID)
	return err
}

// FindTournament finds existing tournament by specified ID
func (pg PostgreSQL) FindTournament(ID uint64, tx *sql.Tx) (*model.Tournament, error) {
	row := tx.QueryRow(queryGetTournament, ID)
	tournament := &model.Tournament{
		ID:      ID,
		Bidders: make([]model.Bidder, 0),
	}
	err := row.Scan(&tournament.IsFinished, &tournament.Deposit)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(queryGetBidders, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		bidder := model.Bidder{
			Backers: make([]string, 0),
		}
		var backersJSON string
		if err := rows.Scan(&bidder.ID, &bidder.Winner, &bidder.Prize, &backersJSON); err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(backersJSON), &bidder.Backers)
		if err != nil {
			return nil, err
		}
		tournament.Bidders = append(tournament.Bidders, bidder)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	pg.Logger.Debugf("Found the tournament %d", tournament.ID)
	return tournament, err
}

// SaveTournament saves a Tournament model
func (pg PostgreSQL) SaveTournament(tournament *model.Tournament, tx *sql.Tx) error {
	stmt, err := tx.Prepare(queryUpdateTournament)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tournament.IsFinished, tournament.Deposit, tournament.ID)
	if err != nil {
		return err
	}
	for _, bidder := range tournament.Bidders {
		backers, err := json.Marshal(bidder.Backers)
		if err != nil {
			return err
		}
		statement, err := tx.Prepare(queryUpsertBidder)
		if err != nil {
			return err
		}
		_, err = statement.Exec(bidder.ID, tournament.ID, bidder.Winner, bidder.Prize, backers)
		if err != nil {
			return err
		}
	}
	pg.Logger.Debugf("Saved the tournament %d", tournament.ID)
	return err
}

// Transaction returns DB transaction control
func (pg PostgreSQL) Transaction() (*sql.Tx, error) {
	return pg.pool.Begin()
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

const (
	// errorCodeDuplicateDatabase is code when create database but it already exists
	errorCodeDuplicateDatabase = "42P04"
	queryGetPlayerBalance      = "SELECT balance FROM players WHERE id = $1"
	queryInsertPlayer          = "INSERT INTO players (id) VALUES ($1)"
	queryUpdatePlayer          = "UPDATE players SET balance = $1 WHERE id = $2"
	queryInsertTournament      = "INSERT INTO tournaments (id) VALUES ($1)"
	queryGetTournament         = "SELECT is_finished, deposit FROM tournaments WHERE id = $1"
	queryUpdateTournament      = "UPDATE tournaments SET is_finished = $1, deposit = $2 WHERE id = $3"
	queryUpsertBidder          = "UPSERT INTO bidders (player_id, tournament_id, winner, prize, backers) VALUES ($1, $2, $3, $4, $5)"
	queryGetBidders            = "SELECT player_id, winner, prize, backers FROM bidders WHERE tournament_id = $1"
)
