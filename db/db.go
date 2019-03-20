package db

import (
	"database/sql"
	"path"
	"sync"
)

// OneSQLite represents the SQLite database to Tramonto One
type OneSQLite struct {
	db  *sql.DB
	mux *sync.Mutex
}

// CreateOneSQLite creates a new instance of the One SQLite database
func CreateOneSQLite(repoPath string) (*OneSQLite, error) {
	// Creates the path to the database
	dbPath := path.Join(repoPath, "database", "one.db")

	// Creates/connects to the database
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Creates instance
	oneSQLite := &OneSQLite{
		db:  conn,
		mux: &sync.Mutex{},
	}

	return oneSQLite, nil
}

// InitTables initializes the tables of the sqlite
func (d *OneSQLite) InitTables() error {
	// Initializes by running migrations
	if err := migrate(d.db); err != nil {
		return err
	}

	return nil
}
