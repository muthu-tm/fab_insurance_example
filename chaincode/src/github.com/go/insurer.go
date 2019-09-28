package main

type Insurer struct {
	InsurerID   string `json:"insurerID"`
	InsurerName string `json:"insurerName"`
	Branch      string `json:"branch,omitempty"`
	State       string `json:"state,omitempty"`
}

// NewInsurer - Creates a new insurer
func NewInsurer(insurerID string, insurerName string) Insurer {
	var insurer Insurer

	insurer.InsurerID = insurerID
	insurer.InsurerName = insurerName

	return insurer
}
