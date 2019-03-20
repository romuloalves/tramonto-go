package db

import (
	"database/sql"

	"github.com/GuiaBolso/darwin"
)

// List of all migrations
var migrations = []darwin.Migration{
	darwin.Migration{
		Version:     1,
		Description: "Create the database",
		Script: `
			create table tests (id integer primary key not null, name text not null, description text not null, ipns text, ipfs text, isOwner integer not null, isFavorite integer not null, createdAt integer not null, updatedAt integer not null);
			create unique index tests_unique_id on tests (id);
		`,
	},
}

// migrate will execute the migrations to the SQLite database
func migrate(db *sql.DB) error {
	// Creates the driver to the SQLite
	driver := darwin.NewGenericDriver(db, darwin.SqliteDialect{})

	// Creates the migration
	d := darwin.New(driver, migrations, nil)

	// Executes the migration
	if err := d.Migrate(); err != nil {
		return err
	}

	return nil
}
