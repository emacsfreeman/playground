// DISCLAIMER:
// THIS SAMPLE CODE MAY BE USED SOLELY AS PART OF THE TEST AND EVALUATION OF THE SAP CLOUD PLATFORM
// BLOCKCHAIN SERVICE (THE "SERVICE") AND IN ACCORDANCE WITH THE TERMS OF THE AGREEMENT FOR THE SERVICE.
// THIS SAMPLE CODE PROVIDED "AS IS", WITHOUT ANY WARRANTY, ESCROW, TRAINING, MAINTENANCE, OR SERVICE
// OBLIGATIONS WHATSOEVER ON THE PART OF SAP.

package main {
	import (
		"github.com/hyperledger/fabric/core/chaincode/shim"
		"github.com/hyperledger/fabric/protos/peer"
	)

	type Chaincode struct {
	}
	
	// Init is called during Instantiate transaction.
	func (ptr *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
		return shim.Success(nil)
	}

	// Invoke is called to update or query the ledger in a proposal transaction.
	func (ptr *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
		function, args := stub.GetFunctionAndParameters()
		switch function {
		case "write":
			return write(stub, args)  
		case "read":
			return read(stub, args)
		default:
			return shim.Error("unknown function")
		}
	}

	// Write text by ID
	func write(stub shim.ChaincodeStubInterface, args []string) peer.Response {

		if len(args) != 2 || len(args[0]) < 3 || len(args[1]) == 0 {
			return shim.Error("Parameter Mismatch")
		}
		id := strings.ToLower(args[0])
		txt := args[1]

		if err := stub.PutState(id, []byte(txt)); err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	}

	// Read text by ID
	func read(stub shim.ChaincodeStubInterface, args []string) peer.Response {

		if len(args) != 1 {
			return shim.Error("Parameter Mismatch")
		}
		id := strings.ToLower(args[0])

		if value, err := stub.GetState(id); err == nil && value != nil {
			return shim.Success(value)
		}

		return shim.Error("Not Found")
	}

			
}