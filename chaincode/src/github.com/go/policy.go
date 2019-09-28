package main

import "reflect"

type Policy struct {
	PolicyID               string `json:"policyID"`
	CustomerID             string `json:"customerID"`
	Type                   string `json:"type"`
	ValidFrom              string `json:"validFrom"`
	ExpiryDate             string `json:"expiryDate"`
	IsCustomerHasOwnership bool   `json:"isCustomerControl"`
}

// NewPolicy - Create a new policy for the customer
func NewPolicy(policyID string, customerID string, customerName string, appliedDate string, expiryDate string, policyType string) Policy {
	var policy Policy

	policy.Type = policyType
	policy.CustomerID = customerID
	policy.PolicyID = policyID
	policy.IsCustomerHasOwnership = false
	policy.ValidFrom = appliedDate
	policy.ExpiryDate = expiryDate

	return policy
}

// IsPolicyEmpty - checks whether the policy struct is empty or not
func (policy Policy) IsPolicyEmpty() bool {
	return reflect.DeepEqual(policy, Policy{})
}
