package main

import (
	"errors"

	shim "github.com/hyperledger/fabric/core/chaincode/shim"
	cid "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
)

func assertAttributeValue(stub shim.ChaincodeStubInterface, attrName string, value string) error {
	_, ok, err := cid.GetAttributeValue(stub, attrName)
	if err != nil {
		return errors.New("There was an error in retrieving the attribute")
	}
	if !ok {
		return errors.New("The identity does not possess the attribute")
	}

	err = cid.AssertAttributeValue(stub, attrName, value)
	if err != nil {
		return err
	}

	return nil
}

func getInitiatorRole(stub shim.ChaincodeStubInterface, attrName string) (string, error) {
	role, ok, err := cid.GetAttributeValue(stub, attrName)
	if err != nil {
		return "", errors.New("There was an error in retrieving the attribute")
	}
	if !ok {
		return "", errors.New("The identity does not possess the attribute")
	}

	return role, nil
}

func checkCustomerOwnership(stub shim.ChaincodeStubInterface, policy Policy) bool {
	return policy.IsCustomerHasOwnership
}
