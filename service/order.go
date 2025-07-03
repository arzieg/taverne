// Package services holds all the services that connect repositories into a business flow
package service

import (
	"taverne/domain/customer"
	"taverne/domain/customer/memory"
)

// OrderConfiguration is an alias for a function that will take in a pointer to an OrderService and modify it
type OrderConfiguration func(os *OrderService) error

// OrderService is a implementation of the OrderService
type OrderService struct {
	customers customer.CustomerRepository
}

// NewOrderService takes a variable amount of OrderConfiguration functions and returns a new OrderService
// Each OrderConfiguration will be called in the order they are passed in
// variadic function (unknown number of parameters)
func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	// create the order service
	os := &OrderService{}
	// apply all configuration passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil

}

// With CustomerRepository applies a given customer repository to the OrderService
func WithCustomerRepository(cr customer.CustomerRepository) OrderConfiguration {
	// return a function that matches the OrderConfiguration aliaas,
	// you need to return this so that the parent function can take in all the needed parameters
	return func(os *OrderService) error {
		os.customers = cr
		return nil
	}
}

// WithMemoryCustomerRepository applies a memory customer repository to the OrderService
func WithMemoryCustomerRepository() OrderConfiguration {
	// create the memory repo, if we needed parameters such as connection strings the could inputet here
	cr := memory.New()
	return WithCustomerRepository(cr)
}
