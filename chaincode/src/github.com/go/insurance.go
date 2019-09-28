/*
* An Insurance chaincode for basci insurance transactions between insurer and customer\
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	shim "github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

// Insurance - chaincode interface to manage Init and Invoke operations
type Insurance struct {
	testMode bool
}

// Init overridden method - it gets called during chaincode instantiation to initialize any data
func (cc *Insurance) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke overriddedn method - it gets called per transaction proposal
func (cc *Insurance) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	var result []byte
	var err error

	if function == CREATENEWPOLICY {
		result, err = newPolicy(stub, args)
	} else if function == QUERYPOLICYBYID {
		result, err = getPolicyByID(stub, args)
	} else if function == QUERYALLPOLICIES {
		result, err = getAllPolicies(stub, args)
	} else if function == GETCUSTOMER {
		result, err = getCustomer(stub, args)
	} else if function == TRANSFERPOLICY {
		result, err = transferPolicy(stub, args)
	} else if function == GETPOLICYPAYMENTS {
		result, err = getPayments(stub, args)
	} else if function == RENEWPOLICY {
		result, err = renewPolicy(stub, args)
	} else if function == "" {
		err = errors.New("Chaincode invoke function name should not be empty")
	} else {
		err = errors.New("Invalid chaincode invoke function name")
	}

	if err != nil {
		fmt.Println("Error occured on chaincode invoke: - ", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(result)
}

func newPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 11 {
		return nil, fmt.Errorf("Incorrect number of arguments for QUERYPOLICYBYID! Expected 11")
	}
	// check whether the transaction proposal initiator is insurer, if no throw error
	err := assertAttributeValue(stub, ROLEATTRIBUTE, INSURERROLE)
	if err != nil {
		return nil, err
	}

	customerID := args[0]
	customerName := args[1]
	insurerID := args[2]
	insurerName := args[3]
	policyID := args[4]
	appliedDate := args[5]
	expiryDate := args[6]
	policyType := args[7]
	amountVal := args[8]
	payType := args[9]
	datePaid := args[10]

	amount, err := strconv.ParseFloat(amountVal, 64)

	if err != nil {
		log.Println("Error occured during amount conversion!!")
		return nil, err
	}

	customer := NewCustomerWithPayment(customerID, customerName, insurerID, policyID, appliedDate, expiryDate, policyType,
		amount, payType, datePaid)

	existingInsurer, insurer, err := IsExistingInsurer(stub, insurerID)
	if err != nil {
		return nil, err
	}
	if existingInsurer {
		policy := NewPolicy(policyID, customerID, customerName, appliedDate, expiryDate, policyType)
		insurer.Policies = append(insurer.Policies, policy)
	} else {
		insurer = NewInsurer(insurerID, insurerName, policyID, customerID, customerName, appliedDate, expiryDate, policyType)
	}

	// convert data into byte array and update the ledger data
	customerAsBytes, _ := json.Marshal(customer)
	insurerAsBytes, _ := json.Marshal(insurer)

	err = PutState(stub, customerID, customerAsBytes)
	if err != nil {
		return nil, err
	}
	err = PutState(stub, insurerID, insurerAsBytes)
	if err != nil {
		return nil, err
	}

	return json.Marshal(customer)
}

func getPolicyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect number of arguments for QUERYPOLICYBYID!! Expected 2")
	}

	insurerID := args[0]
	policyID := args[1]

	policies, err := GetInsurerPolicies(stub, insurerID)
	if err != nil {
		return nil, err
	}

	var policy Policy
	for _, value := range policies {
		if value.PolicyID == policyID {
			policy = value
		}
	}

	isEmpty := policy.IsPolicyEmpty()

	if !isEmpty {
		return json.Marshal(policy)
	}

	return nil, fmt.Errorf("No policy retrieved for the given policyID")
}

func getAllPolicies(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for QUERYALLPOLICIES! Expected 1")
	}

	ID := args[0]

	// get the 'role' attribute value to check the query initator role
	role, err := getInitiatorRole(stub, ROLEATTRIBUTE)
	if err != nil {
		return nil, err
	}

	// check whether the query initiator role is insurer, then get their policies
	if role == INSURERROLE {
		policies, err := GetInsurerPolicies(stub, ID)
		if err != nil {
			return nil, err
		}
		return json.Marshal(policies)
	}

	// check whether the query initiator role is customer, then get their policies
	policies, err := GetCustomerPolicies(stub, ID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(policies)
}

func getCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for GETCUSTOMER! Expected 1")
	}

	customerID := args[0]

	// get the 'role' attribute value to check the query initator role
	role, err := getInitiatorRole(stub, ROLEATTRIBUTE)
	if err != nil {
		return nil, err
	}

	// check the query initiator role in insurer, else send error response as unauthorised
	if role == INSURERROLE {
		customer, err := getState(stub, customerID)
		if err != nil {
			return nil, err
		}
		return customer, nil
	}

	return nil, fmt.Errorf("You are authorised to see this Customer's information")
}

func transferPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("Incorrect number of arguments for TRANSFERPOLICY! Expected 4")
	}

	// check whether the transaction proposal initiator is insurer, if no throw error
	err := assertAttributeValue(stub, ROLEATTRIBUTE, INSURERROLE)
	if err != nil {
		return nil, err
	}

	customerID := args[0]
	insurerID := args[1]
	policyID := args[2]
	transferOwnership := args[3]

	transfer, _ := strconv.ParseBool(transferOwnership)

	// Get customer details from ledger
	customer, err := getCustomerByID(stub, customerID)
	if err != nil {
		return nil, err
	}

	// Get insurer details from ledger
	insurer, err := getInsurerByID(stub, insurerID)
	if err != nil {
		return nil, err
	}

	// Get insurer policies and update the ownership for the given policyID
	for key, value := range insurer.Policies {
		if value.PolicyID == policyID {
			insurer.Policies[key].IsCustomerHasOwnership = transfer
		}
	}

	// Get customer policies and update the ownership for the given policyID
	for key, value := range customer.Policies {
		if value.PolicyID == policyID {
			customer.Policies[key].IsCustomerHasOwnership = transfer
		}
	}

	// convert data into byte array and update the ledger data
	customerAsBytes, _ := json.Marshal(customer)
	insurerAsBytes, _ := json.Marshal(insurer)

	err = PutState(stub, customerID, customerAsBytes)
	if err != nil {
		return nil, err
	}

	err = PutState(stub, insurerID, insurerAsBytes)
	if err != nil {
		return nil, err
	}

	return []byte("Success"), nil
}

func getPayments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for GETPOLICYPAYMENTS! Expected 1")
	}

	customerID := args[0]
	// Get customer details from ledger
	customer, err := getCustomerByID(stub, customerID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(customer.Payments)
}

func renewPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for RENEWPOLICY! Expected 1")
	}
	customerID := args[0]
	policyID := args[1]
	amountVal := args[2]
	payType := args[3]
	datePaid := args[4]
	expiryDate := args[5]

	amount, err := strconv.ParseFloat(amountVal, 64)
	if err != nil {
		log.Println("Error occured during amount conversion!!")
		return nil, err
	}

	// Get customer details from ledger
	customer, err := getCustomerByID(stub, customerID)
	if err != nil {
		return nil, err
	}

	// Set the payment details for the policy
	payment := Payment{PolicyAmount: amount, PaymentType: payType, PolicyID: policyID, PaymentDate: datePaid}
	// Add the pament details to the payments list
	customer.Payments = append(customer.Payments, payment)

	// Get customer policies and update the ownership for the given policyID
	for key, value := range customer.Policies {
		if value.PolicyID == policyID {
			customer.Policies[key].ExpiryDate = expiryDate
		}
	}

	insurerID := customer.InsurerID
	// Get insurer details from ledger
	insurer, err := getInsurerByID(stub, insurerID)
	if err != nil {
		return nil, err
	}

	// Get insurer policies and update the ownership for the given policyID
	for key, value := range insurer.Policies {
		if value.PolicyID == policyID {
			insurer.Policies[key].ExpiryDate = expiryDate
		}
	}

	// convert data into byte array and update the ledger data
	customerAsBytes, _ := json.Marshal(customer)
	insurerAsBytes, _ := json.Marshal(insurer)

	err = PutState(stub, customerID, customerAsBytes)
	if err != nil {
		return nil, err
	}

	err = PutState(stub, insurerID, insurerAsBytes)
	if err != nil {
		return nil, err
	}

	return []byte("Success"), nil
}

func main() {
	if err := shim.Start(new(Insurance)); err != nil {
		fmt.Println("Error creating Insurance ChainCode: ", err)
	} else {
		fmt.Println("Insurance ChainCode was created successfully")
	}

}
