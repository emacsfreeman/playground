package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

// Define the DocRecord Structure, which holds the signature of the document
// signed by issuer, and the time when this record is created
type DocRecord struct {
	Signature   string `json:"signature"`
	Time  string `json:"time"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryDocRecord" {
		return s.queryDocRecord(APIstub, args)
	} else if function == "createDocRecord" {
		return s.createDocRecord(APIstub, args)
	} 

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryDocRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docrecordAsBytes, _ := APIstub.GetState(args[0])

	if docrecordAsBytes == nil {
		return shim.Error("Document not found: " + args[0])
	}

	return shim.Success(docrecordAsBytes)
}

func (s *SmartContract) createDocRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var docrecord = DocRecord{Signature: args[1], Time: time.Now().Format(time.RFC3339)}
	docrecordAsBytes, _ := json.Marshal(docrecord)
	APIstub.PutState(args[0], docrecordAsBytes)

	return shim.Success(nil)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
