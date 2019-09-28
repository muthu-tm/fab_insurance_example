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

func IsExistingInsurer(stub shim.ChaincodeStubInterface, key string) (bool, Insurer, error) {
	bytesData, err := stub.GetState(key)
	var insurer Insurer
	if err != nil {
		return false, insurer, err
	}

	err = json.Unmarshal(bytesData, &insurer)
	if err != nil {
		return false, insurer, err
	}

	return true, insurer, nil
}

func GetInsurerPolicies(stub shim.ChaincodeStubInterface, key string) ([]Policy, error) {
	bytesData, err := stub.GetState(key)
	var insurer Insurer
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytesData, &insurer)
	if err != nil {
		return nil, err
	}

	policies := insurer.Policies
	return policies, nil
}

func GetCustomerPolicies(stub shim.ChaincodeStubInterface, key string) ([]Policy, error) {
	bytesData, err := stub.GetState(key)
	var customer Customer
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytesData, &customer)
	if err != nil {
		return nil, err
	}

	policies := customer.Policies
	return policies, nil
}

func queryState(stub shim.ChaincodeStubInterface, policyID string) ([]byte, error) {
	bytesData, err := stub.GetState(policyID)
	if err != nil {
		return nil, err
	}

	return bytesData, nil
}

func getCustomerByID(stub shim.ChaincodeStubInterface, customerID string) (Customer, error) {
	bytesData, err := getState(stub, customerID)
	var customer Customer
	if err != nil {
		return customer, err
	}

	err = json.Unmarshal(bytesData, &customer)
	if err != nil {
		return customer, err
	}

	return customer, nil
}

func getInsurerByID(stub shim.ChaincodeStubInterface, insurerID string) (Insurer, error) {
	bytesData, err := getState(stub, insurerID)
	var insurer Insurer
	if err != nil {
		return insurer, err
	}

	err = json.Unmarshal(bytesData, &insurer)
	if err != nil {
		return insurer, err
	}

	return insurer, nil
}

func getState(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	resValue, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}

	return resValue, nil
}

/***********************************************************
* getStateByQueryAndDecrypt retrieves the value for the given
* rich query and decrypts it with the supplied entity and
* returns the result of the decryption
***********************************************************/
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
