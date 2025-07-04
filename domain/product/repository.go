// Package product holds the repository and the implementatiopn s for a ProductRepository
package product

import (
	"errors"
	"taverne/aggregate"

	"github.com/google/uuid"
)

var (
	// ErrProductNotFound is returned when a product is not found
	ErrProductNotFound = errors.New("the product was not found")
	// ErrProductAlreadyExist is returned when trying to add a product that already exists
	ErrProductAlreadyExist = errors.New("the product already exists")
)

// ProductRepository is the repository interface to fulfill the use the product aggregate
type ProductRepository interface {
	GetAll() ([]aggregate.Product, error)
	GetByID(id uuid.UUID) (aggregate.Product, error)
	Add(product aggregate.Product) error
	Update(product aggregate.Product) error
	Delete(id uuid.UUID) error
}
