package db

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"

	// Importing to use sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// List of all migrations
var migrations = []darwin.Migration{
	darwin.Migration{
		Version:     1,
		Description: "Create the database",
		Script: `
			CREATE TABLE tests (
    			name             VARCHAR (6) NOT NULL,
    			description      TEXT        NOT NULL,
    			secret           TEXT        NOT NULL,
    			ipfs_hash        VARCHAR,
    			ipns_hash        VARCHAR,
    			is_key_generated BOOLEAN     NOT NULL
                                 			DEFAULT false,
    			is_owner         BOOLEAN     NOT NULL
                                 			DEFAULT true,
    			is_favorite      BOOLEAN     NOT NULL
                                 			DEFAULT false,
    			created_at       TIMESTAMP   NOT NULL
                                 			DEFAULT (CURRENT_TIMESTAMP),
    			updated_at       TIMESTAMP   NOT NULL
                                 			DEFAULT (CURRENT_TIMESTAMP),
    			is_active        BOOLEAN     NOT NULL
                                 			DEFAULT true
			);
		`,
	},
}

// migrate will execute the migrations to the SQLite database
func migrate(db *sqlx.DB) error {
	// Creates the driver to the SQLite
	driver := darwin.NewGenericDriver(db.DB, darwin.SqliteDialect{})

	// Creates the migration
	d := darwin.New(driver, migrations, nil)

	// Executes the migration
	if err := d.Migrate(); err != nil {
		return err
	}

	return nil
}
