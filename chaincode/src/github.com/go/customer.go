package main

type Customer struct {
	CustomerID string `json:"customerID"`
	Name       string `json:"firstname"`
	Contact    string `json:"contact,omitempty"`
}

// NewCustomer - Creates a new customer
func NewCustomer(customerID string, customerName string) Customer {
	var customer Customer

	customer.CustomerID = customerID
	customer.Name = customerName
	return customer
}
