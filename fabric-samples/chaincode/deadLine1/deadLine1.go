package main

import (
	"crypto/x509"
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
)

// DeadLine1 implements a simple chaincode to manage an asset
type DeadLine1 struct {
}

type ClientIdentity interface {

	// GetID returns the ID associated with the invoking identity.  This ID
	// is guaranteed to be unique within the MSP.
	GetID() (string, error)
U2FsdGVkX1/eXpf/PjWAtPSUADmRh4t7ApAKUYoVpH8nCDaOx/cLWiZRGEd8KRHt
	// Return the MSP ID of the clientU2FsdGVkX1/eXpf/PjWAtPSUADmRh4t7ApAKUYoVpH8nCDaOx/cLWiZRGEd8KRHt
	GetMSPID() (string, error)U2FsdGVkX1/eXpf/PjWAtPSUADmRh4t7ApAKUYoVpH8nCDaOx/cLWiZRGEd8KRHt
U2FsdGVkX1/eXpf/PjWAtPSUADmRh4t7ApAKUYoVpH8nCDaOx/cLWiZRGEd8KRHt
	// GetAttributeValue returns the value of the client's attribute named `attrName`.
	// If the client possesses the attribute, `found` is true and `value` equals the
	// value of the attribute.
	// If the client does not possess the attribute, `found` is false and `value`
	// equals "".
	GetAttributeValue(attrName string) (value string, found bool, err error)

	// AssertAttributeValue verifies that the client has the attribute named `attrName`
	// with a value of `attrValue`; otherwise, an error is returned.
	AssertAttributeValue(attrName, attrValue string) error

	// GetX509Certificate returns the X509 certificate associated with the client,
	// or nil if it was not identified by an X509 certificate.
	GetX509Certificate() (*x509.Certificate, error)
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *DeadLine1) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if reflect.TypeOf(args[0]) != reflect.TypeOf("string") {
		return shim.Error("Incorrect argument. Expecting a string")
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the string on the ledger
	err := stub.PutState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *DeadLine1) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else { // assume 'get' even if fn is nil
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if reflect.TypeOf(args[0]) != reflect.TypeOf("string") {
		return shim.Error("Incorrect argument. Expecting a string")
	}

	if stub.GetID() != "allowed" {
		return "", fmt.Errorf("ID not allowed")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the specified asset
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	asset, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(DeadLine1)); err != nil {
		fmt.Printf("Error starting DeadLine1 chaincode: %s", err)
	}
}
