package db

import (
	"path"
	"sync"

	"github.com/jmoiron/sqlx"
)

// OneSQLite represents the SQLite database to Tramonto One
type OneSQLite struct {
	dbPath string
	db     *sqlx.DB
	mux    *sync.Mutex
}

// OpenOneSQLite creates a new instance of the One SQLite database
func OpenOneSQLite(repoPath string) (*OneSQLite, error) {
	// Creates the path to the database
	dbPath := path.Join(repoPath, "one.db")

	// Creates/connects to the database
	conn, err := sqlx.Open("sqlite3", dbPath)
	// conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Creates instance
	oneSQLite := &OneSQLite{
		dbPath: dbPath,
		db:     conn,
		mux:    new(sync.Mutex),
	}

	return oneSQLite, nil
}

// MigrateTables initializes and/or migrates the tables of the sqlite
func (d *OneSQLite) MigrateTables() error {
	d.mux.Lock()
	defer d.mux.Unlock()

	// Initializes by running migrations
	if err := migrate(d.db); err != nil {
		return err
	}

	return nil
}
