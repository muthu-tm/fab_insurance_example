/*
* An Insurance chaincode for basic policy transactions between insurer and customer
 */
package main

import (
	"bytes"
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

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (cc *Insurance) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
func (cc *Insurance) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// GetFunctionAndParameters returns the first argument as the function name
	// and the rest of the arguments as parameters in a string array.
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
	} else if function == QUERYCUSTOMERPOLICY {
		result, err = getCustomerPolicy(stub, args)
	} else if function == QUERYPOLICYHISTORY {
		result, err = getPolicyHistory(stub, args)
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

	existingPolicy, err := IsExistingPolicy(stub, policyID)

	if existingPolicy {
		return nil, fmt.Errorf("Found a policy for the given policyID")
	}
	policy := NewPolicyWithPayment(customerID, customerName, insurerID, insurerName, policyID, appliedDate, expiryDate, policyType,
		amount, payType, datePaid)

	// convert data into byte array and update the ledger data
	policyAsBytes, _ := json.Marshal(policy)

	err = PutState(stub, policyID, policyAsBytes)
	if err != nil {
		return nil, err
	}

	return json.Marshal(policy)
}

func getPolicyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for QUERYPOLICYBYID!! Expected 2")
	}

	// check whether the transaction proposal initiator is insurer, if no throw error
	err := assertAttributeValue(stub, ROLEATTRIBUTE, INSURERROLE)
	if err != nil {
		return nil, err
	}

	policyID := args[0]

	policy, err := queryState(stub, policyID)
	if err != nil {
		log.Println("No policy retrieved for the given policyID")
		return nil, err
	}

	return json.Marshal(policy)
}

func getAllPolicies(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for QUERYALLPOLICIES! Expected 1")
	}

	// check whether the transaction proposal initiator is insurer, if no throw error
	err := assertAttributeValue(stub, ROLEATTRIBUTE, INSURERROLE)
	if err != nil {
		return nil, err
	}

	insurerID := args[0]
	query := fmt.Sprintf("{\"selector\" : {\"insurer\": { \"insurerID\": \"%s\"}}}", insurerID)

	return getStateByQuery(stub, query)
}

func getCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for GETCUSTOMER! Expected 1")
	}

	// check whether the transaction proposal initiator is insurer, if no throw error
	err := assertAttributeValue(stub, ROLEATTRIBUTE, INSURERROLE)
	if err != nil {
		log.Println("Attribute assertion failure! Not authorised to see this Customer's information")
		return nil, err
	}
	customerID := args[0]
	query := fmt.Sprintf("{\"selector\" : {\"customer\": { \"customerID\": \"%s\"}}}", customerID)

	return getStateByQuery(stub, query)
}

func transferPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect number of arguments for TRANSFERPOLICY! Expected 2")
	}

	// check whether the transaction proposal initiator is insurer, if no throw error
	err := assertAttributeValue(stub, ROLEATTRIBUTE, INSURERROLE)
	if err != nil {
		return nil, err
	}

	policyID := args[0]
	transferOwnership := args[1]

	transfer, _ := strconv.ParseBool(transferOwnership)

	// Get policy details from ledger
	policy, err := GetPolicyByID(stub, policyID)
	if err != nil {
		return nil, err
	}
	policy.IsCustomerHasOwnership = transfer

	// convert data into byte array and update the ledger data
	byteValue, _ := json.Marshal(policy)

	err = PutState(stub, policyID, byteValue)
	if err != nil {
		return nil, err
	}

	return json.Marshal(policy)
}

func getPayments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for GETPOLICYPAYMENTS! Expected 1")
	}

	// check whether the transaction proposal initiator is customer, if no throw error
	err := assertAttributeValue(stub, ROLEATTRIBUTE, CUSTOMERROLE)
	if err != nil {
		return nil, err
	}

	policyID := args[0]
	// Get policy details from ledger
	policy, err := GetPolicyByID(stub, policyID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(policy.Payments)
}

func renewPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("Incorrect number of arguments for RENEWPOLICY! Expected 5")
	}
	policyID := args[0]
	amountVal := args[1]
	payType := args[2]
	datePaid := args[3]
	expiryDate := args[4]

	amount, err := strconv.ParseFloat(amountVal, 64)
	if err != nil {
		log.Println("Error occured during amount conversion!!")
		return nil, err
	}

	// Get policy details from ledger
	policy, err := GetPolicyByID(stub, policyID)
	if err != nil {
		return nil, err
	}

	// Set the payment details for the policy
	payment := Payment{PolicyAmount: amount, PaymentType: payType, PolicyID: policyID, PaymentDate: datePaid}
	// Add the pament details to the payments list and update the expiry date
	policy.Payments = append(policy.Payments, payment)
	policy.ExpiryDate = expiryDate

	// convert data into byte array and update the ledger data
	policyAsBytes, _ := json.Marshal(policy)

	err = PutState(stub, policyID, policyAsBytes)
	if err != nil {
		return nil, err
	}

	return json.Marshal(policy)
}

func getCustomerPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for QUERYCUSTOMERPOLICY! Expected 1")
	}
	customerID := args[0]

	query := fmt.Sprintf("{\"selector\" : {\"customer\": { \"customerID\": \"%s\"}}}", customerID)

	return getStateByQuery(stub, query)
}

func getPolicyHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments for QUERYPOLICYHISTORY! Expected 1")
	}
	policyID := args[0]

	// GetHistoryForKey returns a history of key values across time.
	resultsIterator, err := stub.GetHistoryForKey(policyID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing the retrieved History Records
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(policyID)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Value\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}

func main() {
	if err := shim.Start(new(Insurance)); err != nil {
		fmt.Println("Error creating Insurance ChainCode: ", err)
	} else {
		fmt.Println("Insurance ChainCode was created successfully")
	}

}
