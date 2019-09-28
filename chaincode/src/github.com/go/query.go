package main

import (
	"bytes"
	"encoding/json"
	"log"

	shim "github.com/hyperledger/fabric/core/chaincode/shim"
)

func PutState(stub shim.ChaincodeStubInterface, key string, value []byte) error {

	// TODO: Encrypt the value before storing into the ledger using bccsp
	// https://godoc.org/github.com/hyperledger/fabric/bccsp

	return stub.PutState(key, value)
}

// IsExistingPolicy - checks whether the policy already exists
func IsExistingPolicy(stub shim.ChaincodeStubInterface, key string) (bool, error) {
	// If the key does not exist in the state database, (nil, nil) is returned.
	byteVal, err := stub.GetState(key)
	if err != nil {
		return false, err
	}
	if byteVal != nil {
		return true, nil
	}

	return false, nil
}

func queryState(stub shim.ChaincodeStubInterface, policyID string) ([]byte, error) {
	bytesData, err := stub.GetState(policyID)
	if err != nil {
		return nil, err
	}

	return bytesData, nil
}

// GetPolicyByID - Get the policy and returns the Policy struct after conversion
func GetPolicyByID(stub shim.ChaincodeStubInterface, policyID string) (Policy, error) {
	bytesData, err := getState(stub, policyID)
	var policy Policy
	if err != nil {
		return policy, err
	}

	err = json.Unmarshal(bytesData, &policy)
	if err != nil {
		return policy, err
	}

	return policy, nil
}

func getState(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	resValue, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}

	return resValue, nil
}

// getStateByQueryAndDecrypt retrieves the value for the given rich query
func getStateByQuery(stub shim.ChaincodeStubInterface, query string) ([]byte, error) {

	log.Println("\n >>  QueryString:- ", query)

	resultsIterator, err := stub.GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing the retrieved QueryRecords
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
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Value\":")

		// The retrieved value record is a JSON object, so we write as-is
		var policy Policy
		resValue, err := json.Marshal(policy)
		buffer.WriteString(string(resValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}
