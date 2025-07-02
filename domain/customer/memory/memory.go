package memory

import (
	"fmt"
	"sync"
	"taverne/aggregate"
	"taverne/domain/customer"

	"github.com/google/uuid"
)

// Let’s first set up the correct structure in the memory file, we want to create a struct that has methods
// to fulfill the CustomerRepository, and let’s not forget the factory to create a new repository.

// MemoryRepository fulfills the CustomerRepository interface
type MemoryRepository struct {
	customers map[uuid.UUID]aggregate.Customer
	sync.Mutex
}

// new is a factory function to generate a new repository of customers
func New() *MemoryRepository {
	return &MemoryRepository{
		customers: make(map[uuid.UUID]aggregate.Customer),
	}
}

// Get finds a customer by ID
func (mr *MemoryRepository) Get(id uuid.UUID) (aggregate.Customer, error) {
	if customer, ok := mr.customers[id]; ok {
		return customer, nil
	}
	return aggregate.Customer{}, customer.ErrCustomerNotFound
}

// Add will add a new customer to the repository
func (mr *MemoryRepository) Add(c aggregate.Customer) error {
	if mr.customers == nil {
		// saftey check if customers is not create
		mr.Lock()
		mr.customers = make(map[uuid.UUID]aggregate.Customer)
		mr.Unlock()
	}
	// Make sure customer isn't already in the repository
	if _, ok := mr.customers[c.GetID()]; ok {
		return fmt.Errorf("customer already exists: %w", customer.ErrFailedToAddCustomer)
	}
	mr.Lock()
	mr.customers[c.GetID()] = c
	mr.Unlock()
	return nil
}

// Update will replace an existing customer information with the new customer information
func (mr *MemoryRepository) Update(c aggregate.Customer) error {
	// Make sure Customer is in the repository
	if _, ok := mr.customers[c.GetID()]; !ok {
		return fmt.Errorf("customer does not exists: %w", customer.ErrUpdateCustomer)
	}
	mr.Lock()
	mr.customers[c.GetID()] = c
	mr.Unlock()
	return nil
}
