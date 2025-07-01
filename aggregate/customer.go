package aggregate

import (
	"taverne/entity"
	"taverne/valueobject"
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
