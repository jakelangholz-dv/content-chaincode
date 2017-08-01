package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ContentChainCode initial struct
type ContentChainCode struct {
}

func main() {
	err := shim.Start(new(ContentChainCode))
	if err != nil {
		fmt.Printf("Could not initialize ContentChainCode: %s", err)
	}
}

// Init creates the initial state of the ContentChainCode
func (t *ContentChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) > 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	return nil, nil
}

// Invoke is the entry point for all state changing functions
func (t *ContentChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "create" {
		return t.create(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	}
	fmt.Println("Invalid function request received: " + function)

	return nil, errors.New("Unable to invoke function: " + function)
}

// Query is our entry point for reads from the chaincode state
func (t *ContentChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// create inserts a new content and owner into the chaincode state
func (t *ContentChainCode) create(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var bookID, ownerID string
	var err error
	fmt.Println("Creating a new content entry!")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2: bookID and ownerID")
	}

	bookID = args[0]
	ownerID = args[1]
	err = stub.PutState(bookID, []byte(ownerID))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// transfer changes the state of a content entry in the chaincode state
func (t *ContentChainCode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var bookID, oldOwnerID, newOwnerID string
	var err error
	fmt.Println("running write()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3: bookID, oldOwnerID, newOwnerID")
	}

	bookID = args[0] //rename for funsies
	oldOwnerID = args[1]
	newOwnerID = args[2]
	oldOwnerFromState, err := stub.GetState(bookID)
	if oldOwnerID == string(oldOwnerFromState) {
		err = stub.PutState(bookID, []byte(newOwnerID))
	} else {
		err = errors.New("Parameter oldOwnerID does not match current state")
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read queries the chaincode state for the current owner of content
func (t *ContentChainCode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var bookID, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting bookID to query")
	}

	bookID = args[0]
	valAsbytes, err := stub.GetState(bookID)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + bookID + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
