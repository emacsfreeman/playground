package main

import (
	"fmt"

	fc "github.com/chaincode/fabcrypt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "testEncrypt" {
		return s.testEncrypt(APIstub)
	} else if function == "testDecrypt" {
		return s.testDecrypt(APIstub)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) testEncrypt(APIstub shim.ChaincodeStubInterface) sc.Response {
	value := "WakandaForeva"
	key := "testId"
	err := fc.Encrypter(APIstub, key, []byte(value))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) testDecrypt(APIstub shim.ChaincodeStubInterface) sc.Response {
	key := "testId"
	val, err := fc.Decrypter(APIstub, key)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(val)
}
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting SmartContract chaincode: %s", err)
	}
}
