package main

type Customer struct {
	CustomerID string    `json:"customerID"`
	Name       string    `json:"firstname"`
	Contact    string    `json:"contact,omitempty"`
	Payments   []Payment `json:"payments,omitempty"`
	Policies   []Policy  `json:"policies,omitempty"`
	InsurerID  string    `json:"insurer"`
}

type Payment struct {
	PolicyAmount float64 `json:"amount"`
	PaymentType  string  `json:"paymenttype"`
	PolicyID     string  `json:"policyID"`
	PaymentDate  string  `json:"datePaid"`
}

// newCustomer - Create a new customer without payment
func newCustomer(customerID string, customerName string, insurerID string, policyID string, appliedDate string, expiryDate string, policyType string) Customer {
	var customer Customer

	customer.CustomerID = customerID
	customer.Name = customerName
	customer.InsurerID = insurerID

	policy := NewPolicy(policyID, customerID, customerName, appliedDate, expiryDate, policyType)

	customer.Policies = append(customer.Policies, policy)

	return customer
}

// NewCustomerWithPayment - Create a new customer with policy payment
func NewCustomerWithPayment(customerID string, customerName string, insurerID string, policyID string, appliedDate string, expiryDate string, policyType string, amount float64, payType string, datePaid string) Customer {
	// Create a new customer
	customer := newCustomer(customerID, customerName, insurerID, policyID, appliedDate, expiryDate, policyType)

	// Set the payment details for the policy
	payment := Payment{PolicyAmount: amount, PaymentType: payType, PolicyID: policyID, PaymentDate: datePaid}

	// Add the pament details to the payments list
	customer.Payments = append(customer.Payments, payment)

	return customer
}
