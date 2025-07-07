package aggregate

import (
	"errors"
	"taverne/entity"
	"taverne/valueobject"

	"github.com/google/uuid"
)

var (
	// ErrInvalidPerson is returned when the person is not valid in the NewCustomer factory
	ErrInvalidPerson = errors.New("a customer has to have an valid person")
)

type Customer struct {
	// person is the root entiry of a customer
	// wich meachs the person.ID is the main identifier for this aggregation
	person *entity.Person
	// a customer can hold many products
	products []*entity.Item
	// a customer can perform many transactions
	transactions []valueobject.Transaction
}

// NewCustomer is a factory to create a new Customer aggregate
// It will validate that the name is not empty
func NewCustomer(name string) (Customer, error) {
	if name == "" {
		return Customer{}, ErrInvalidPerson
	}

	person := &entity.Person{
		Name: name,
		ID:   uuid.New(),
	}

	// Create a customer object and initialize all the values to avoid nil pointer exceptions
	return Customer{
		person:       person,
		products:     make([]*entity.Item, 0),
		transactions: make([]valueobject.Transaction, 0),
	}, nil
}

// GetID returns the customers root entity ID
func (c *Customer) GetID() uuid.UUID {
	return c.person.ID
}

// SetID setzs the root ID
func (c *Customer) SetID(id uuid.UUID) {
	if c.person == nil {
		c.person = &entity.Person{}
	}
	c.person.ID = id
}

// SetName changes the name of the customer
func (c *Customer) SetName(name string) {
	if c.person == nil {
		c.person = &entity.Person{}
	}
	c.person.Name = name
}

// SetName changes the name of the customer
func (c *Customer) GetName() string {
	return c.person.Name
}
