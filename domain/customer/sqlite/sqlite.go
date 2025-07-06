package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"taverne/aggregate"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
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
	db, err := sql.Open("sqlite3", "/tmp/ddd.db")
	//defer db.Close()
	if err != nil {
		return nil, err
	}

	// create tabke customers
	_, err = db.ExecContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS customer (
			id TEXT PRIMARY KEY, 
			name TEXT NOT NULL,
			age INT
		)`,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating table customer, got %v", err)
	}

	return &SqliteRepository{
		db: db,
	}, nil

}

func (sr *SqliteRepository) Get(id uuid.UUID) (aggregate.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT id, name FROM customer WHERE id = ?`
	var result sqliteCustomer

	err := sr.db.QueryRowContext(ctx, query, id).Scan(&result.ID, &result.Name)
	if err != nil {
		return aggregate.Customer{}, err
	}
	return result.ToAggregate(), nil
}

func (sr *SqliteRepository) Add(c aggregate.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	internal := NewFromCustomer(c)
	query := `INSERT INTO customer (id, name, age) VALUES (?, ?, ?)`
	_, err := sr.db.ExecContext(ctx, query, internal.ID.String(), internal.Name, nil)
	if err != nil {
		return fmt.Errorf("insert into customers failed, got %v\n", err)
	}
	return nil
}

func (sr *SqliteRepository) Update(c aggregate.Customer) error {
	panic("to implement")
}
