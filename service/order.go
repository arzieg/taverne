// Package services holds all the services that connect repositories into a business flow
package service

import (
	"context"
	"log"
	"taverne/aggregate"
	"taverne/domain/customer"
	"taverne/domain/customer/memory"
	"taverne/domain/customer/sqlite"
	"taverne/domain/product"
	prodmemory "taverne/domain/product/memory"

	"github.com/google/uuid"
)

// OrderConfiguration is an alias for a function that will take in a pointer to an OrderService and modify it
type OrderConfiguration func(os *OrderService) error

// OrderService is a implementation of the OrderService
type OrderService struct {
	customers customer.CustomerRepository
	products  product.ProductRepository
}

// NewOrderService takes a variable amount of OrderConfiguration functions and returns a new OrderService
// Each OrderConfiguration will be called in the order they are passed in.
// It is a variadic function (aka unknown number of parameters)
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

// WithMemotyProductRepository adds a in memory product repo and adds all input products
func WithMemoryProductRepository(products []aggregate.Product) OrderConfiguration {
	return func(os *OrderService) error {
		// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
		pr := prodmemory.New()

		// Add Items to repo
		for _, p := range products {
			err := pr.Add(p)
			if err != nil {
				return err
			}
		}
		os.products = pr
		return nil
	}
}

func WithSQLiteCustomerRepository(connectionString string) OrderConfiguration {
	return func(os *OrderService) error {
		// Create the sqlite repo, if we needed parameters, such as connection strings they could be inputted here
		cr, err := sqlite.New(context.Background(), connectionString)
		if err != nil {
			return err
		}
		os.customers = cr
		return nil
	}
}

// CreateOrder will chaintogether all repositories to create a order for a customer
// will return the collected price of all Products
func (o *OrderService) CreateOrder(customerID uuid.UUID, productIDs []uuid.UUID) (float64, error) {
	// Get the customer
	c, err := o.customers.Get(customerID)
	if err != nil {
		return 0, err
	}

	// Get each product
	var products []aggregate.Product
	var price float64
	for _, id := range productIDs {
		p, err := o.products.GetByID(id)
		if err != nil {
			return 0, err
		}
		products = append(products, p)
		price += p.GetPrice()
	}

	// All Products exists in store, now we can create the order
	log.Printf("Customer: %s has ordered %d products", c.GetID(), len(products))

	return price, nil
}
