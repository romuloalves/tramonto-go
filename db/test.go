package db

import (
	"errors"
	"time"

	// Importing to use sqlite3
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/tramonto-one/go-tramonto/entities"
)

type dbTest struct {
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	Secret         string    `db:"secret"`
	IpfsHash       string    `db:"ipfs_hash"`
	IpnsHash       string    `db:"ipns_hash"`
	IsKeyGenerated bool      `db:"is_key_generated"`
	IsOwner        bool      `db:"is_owner"`
	IsFavorite     bool      `db:"is_favorite"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	IsActive       bool      `db:"is_active"`
}

// InsertTest inserts a new test to the database
func (db *OneSQLite) InsertTest(test entities.Test) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Starts transaction
	tx := db.db.MustBegin()

	// Inserts data
	sqlResult := tx.MustExec(`
		INSERT INTO tests (name, description, secret, ipfs_hash, ipns_hash, is_key_generated)
		VALUES ($1, $2, $3, $4, $5, $6);`,
		test.Metadata.Name, test.Metadata.Description, test.Secret, test.Ipfs, test.Ipns, test.IpnsKeyCreated)

	_, err := sqlResult.RowsAffected()
	if err != nil {
		panic(err)
	}

	// Commits
	if err := tx.Commit(); err != nil {
		return errors.New("Error inserting test: " + err.Error())
	}

	return nil
}

// FindTests finds all active tests
func (db *OneSQLite) FindTests() ([]entities.Test, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Array to receive the select results
	tests := []dbTest{}

	// Executes the select
	if err := db.db.Select(&tests, `
		SELECT *
		FROM tests
		WHERE is_active = 1
		ORDER BY updated_at DESC`); err != nil {
		return []entities.Test{}, errors.New("Error finding tests: " + err.Error())
	}

	// Parses to entity
	result := []entities.Test{}

	for _, test := range tests {
		result = append(result, entities.Test{
			Ipfs:           test.IpfsHash,
			Ipns:           test.IpnsHash,
			IpnsKeyCreated: test.IsKeyGenerated,
			Secret:         test.Secret,
			Metadata: entities.Metadata{
				Name:        test.Name,
				Description: test.Description,
			},
		})
	}

	return result, nil
}
