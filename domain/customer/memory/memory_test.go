package memory

import (
	"taverne/aggregate"
	"taverne/domain/customer"
	"testing"

	"github.com/google/uuid"
)

func TestMemory_GetCustomer(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	cust, err := aggregate.NewCustomer("Donald")
	if err != nil {
		t.Fatal(err)
	}
	id := cust.GetID()

	// create the repo to use and add some test data to it for testing
	// Skip Factory for this
	repo := MemoryRepository{
		customers: map[uuid.UUID]aggregate.Customer{
			id: cust,
		},
	}
	testCases := []testCase{
		{
			name:        "No customer by ID",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: customer.ErrCustomerNotFound,
		}, {
			name:        "Customer By ID",
			id:          id,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := repo.Get(tc.id)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

		})
	}

}

func TestMemory_AddCustomer(t *testing.T) {
	type testCase struct {
		name        string
		cust        string
		expectedErr error
	}

	testCases := []testCase{
		{
			name:        "Add Customer",
			cust:        "Donald",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := MemoryRepository{
				customers: map[uuid.UUID]aggregate.Customer{},
			}

			cust, err := aggregate.NewCustomer(tc.cust)
			if err != nil {
				t.Fatal(err)
			}

			err = repo.Add(cust)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

			found, err := repo.Get(cust.GetID())
			if err != nil {
				t.Fatal(err)
			}
			if found.GetID() != cust.GetID() {
				t.Errorf("Expected %v, got %v", cust.GetID(), found.GetID())
			}
		})
	}
}
