package main

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

const DECKEY = "DECKEY"
const ENCKEY = "ENCKEY"
const IV = "IV"
const KEY = "key"
const VALUE = "value"

// Invoke for this chaincode exposes functions to ENCRYPT, DECRYPT transactional
// data. The Initialization Vector (IV) is passed in as a parm to
// ensure peers have deterministic data.
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// get arguments and transient
	f, _ := stub.GetFunctionAndParameters()
	// for the arguments correlated to encryption,
	// they are supposed to be store in transient map instead of arguments array.
	// please remind that data in argument array persist in the ledger.
	tMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve transient, err %s", err))
	}

	switch f {

	case "ENCRYPT":
		// make sure there's a key in transient - the assumption is that
		// it's associated to the string "ENCKEY"
		if _, in := tMap[ENCKEY]; !in {
			return shim.Error(fmt.Sprintf("Expected transient encryption key %s", ENCKEY))
		}

		// We need to decode the encryption key from base64 string.
		// Usually, encryption key is encoded in base64 string before sending in a REST request.
		encKey, err := base64.StdEncoding.DecodeString(string(tMap[ENCKEY]))
		if err != nil {
			return shim.Error(err.Error())
		}

		// Also for the IV.
		iv, err := base64.StdEncoding.DecodeString(string(tMap[IV]))
		if err != nil {
			return shim.Error(err.Error())
		}

		if _, ok := tMap[KEY]; !ok {
			return shim.Error("cannot find state key")
		}
		if _, ok := tMap[VALUE]; !ok {
			return shim.Error("cannot find state value")
		}

		args := []string{string(tMap[KEY]), string(tMap[VALUE])}

		return s.Encrypt(stub, args[0:], encKey, iv)

	case "DECRYPT":
		// make sure there's a key in transient - the assumption is that
		// it's associated to the string "DECKEY"
		if _, in := tMap[DECKEY]; !in {
			return shim.Error(fmt.Sprintf("Expected transient decryption key %s", DECKEY))
		}

		// the decryption key is similar to encryption key which is wrapped in base64 string.
		decKey, err := base64.StdEncoding.DecodeString(string(tMap[DECKEY]))
		if err != nil {
			return shim.Error(err.Error())
		}

		iv, err := base64.StdEncoding.DecodeString(string(tMap[IV]))
		if err != nil {
			return shim.Error(err.Error())
		}

		if _, ok := tMap[KEY]; !ok {
			return shim.Error("cannot find state key")
		}

		args := []string{string(tMap[KEY])}
		return s.Decrypt(stub, args[0:], decKey, iv)

	default:
		return shim.Error(fmt.Sprintf("Unsupported function %s", f))
	}
}

func (s *SmartContract) Encrypt(stub shim.ChaincodeStubInterface, args []string, encKey, IV []byte) peer.Response {

	if len(args) != 2 {
		return shim.Error("Expected 2 parameters to function Encrypt")
	}
	key := args[0]
	value := args[1]

	if err := Encrypt(stub, key, []byte(value), encKey, IV); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) Decrypt(stub shim.ChaincodeStubInterface, args []string, decKey, IV []byte) peer.Response {

	if len(args) != 1 {
		return shim.Error("Expected 1 parameters to function Decrypt")
	}

	key := args[0]

	if val, err := Decrypt(stub, key, decKey, IV); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(val)
	}
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting SmartContract chaincode: %s", err)
	}
}
