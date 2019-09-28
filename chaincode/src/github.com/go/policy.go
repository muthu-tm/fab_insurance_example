package main

type Policy struct {
	PolicyID               string    `json:"policyID"`
	Type                   string    `json:"type"`
	ValidFrom              string    `json:"validFrom"`
	ExpiryDate             string    `json:"expiryDate"`
	Customer               Customer  `json:"customer"`
	Insurer                Insurer   `json:"insurer"`
	Payments               []Payment `json:"payments,omitempty"`
	IsCustomerHasOwnership bool      `json:"isCustomerControl"`
}

type Payment struct {
	PolicyAmount float64 `json:"amount"`
	PaymentType  string  `json:"paymenttype"`
	PolicyID     string  `json:"policyID"`
	PaymentDate  string  `json:"datePaid"`
}

// NewPolicy - Create a new policy
func NewPolicy(policyID string, customerID string, customerName string, appliedDate string, expiryDate string, policyType string, insurerID string, insurerName string) Policy {
	var policy Policy

	policy.Type = policyType
	policy.PolicyID = policyID
	policy.IsCustomerHasOwnership = false
	policy.ValidFrom = appliedDate
	policy.ExpiryDate = expiryDate

	policy.Customer = NewCustomer(customerID, customerName)
	policy.Insurer = NewInsurer(insurerID, insurerName)

	return policy
}

// NewPolicyWithPayment - Create a new policy with payment
func NewPolicyWithPayment(customerID string, customerName string, insurerID string, insurerName string, policyID string, appliedDate string, expiryDate string, policyType string, amount float64, payType string, datePaid string) Policy {

	policy := NewPolicy(policyID, customerID, customerName, appliedDate, expiryDate, policyType, insurerID, insurerName)

	// Set the payment details for the policy
	payment := Payment{PolicyAmount: amount, PaymentType: payType, PolicyID: policyID, PaymentDate: datePaid}

	// Add the pament details to the payments list
	policy.Payments = append(policy.Payments, payment)

	return policy
}
