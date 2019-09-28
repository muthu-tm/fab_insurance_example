package main

type Insurer struct {
	InsurerID   string   `json:"insurerID"`
	InsurerName string   `json:"insurerName"`
	Policies    []Policy `json:"policies,omitempty"`
	Branch      string   `json:"branch,omitempty"`
	State       string   `json:"state,omitempty"`
}

// NewInsurer - Creates new insurer
func NewInsurer(insurerID string, insurerName string, policyID string, customerID string, customerName string, appliedDate string, expiryDate string, policyType string) Insurer {
	var insurer Insurer

	insurer.InsurerID = insurerID
	insurer.InsurerName = insurerName

	policy := NewPolicy(policyID, customerID, customerName, appliedDate, expiryDate, policyType)

	insurer.Policies = append(insurer.Policies, policy)

	return insurer
}
