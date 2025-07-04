package sqlite

import (
	"context"
	"database/sql"
	"taverne/aggregate"

	"github.com/gofrs/uuid"
)

type SqliteRepository struct {
	db *sql.DB
}

// sqliteCustomer is an internal type that is used to store a CustomerAggregate
// we make an internal struct for this to avoid coupling this sqlite implementation to the customeraggregate.
// sqlite uses
type sqliteCustomer struct {
	ID   uuid.UUID
	Name string
}

// NewFromCustomer takes in a aggregate and converts into internal structure
func NewFromCustomer(c aggregate.Customer) sqliteCustomer {
	return sqliteCustomer{
		ID:   c.GetID(),
		Name: c.GetName(),
	}
}

// ToAggregate converts into a aggregate.Customer
// ths could validate all values present etc
func (s sqliteCustomer) ToAggregate() aggregate.Customer {
	c := aggregate.Customer{}

	c.SetID(s.ID)
	c.SetName(s.Name)

	return c
}

// Create a new sqlite repository
func New(ctx context.Context, connectionString string) (*SqliteRepository, error) {
	db, err := sql.Open("sqlite3", "ddd")
	if err != nil {
		return nil, err
	}

	customers, err := db.ExecContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS customers (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			title TEXT NOT NULL, 
			artist TEXT NOT NULL, 
			price REAL NOT NULL
		)`,
	)

}

func initDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS album (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			title TEXT NOT NULL, 
			artist TEXT NOT NULL, 
			price REAL NOT NULL
		)`,
	)
	if err != nil {
		return err
	}
	return nil
}
