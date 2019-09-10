package main

import (
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
)

// EncCC example simple Chaincode implementation of a chaincode that uses encryption/signatures
var (
	bccspInst bccsp.BCCSP
)

func init() {
	factory.InitFactories(nil)
	bccspInst = factory.GetDefault()
}

// Encrypt exposes how to write state to the ledger after having
// encrypted it with an AES 256 bit key that has been provided to the chaincode through the
// transient field
func Encrypt(stub shim.ChaincodeStubInterface, key string, cleartextValueAsByte []byte, encKey, IV []byte) error {
	// create the encrypter entity - we give it an ID, the bccsp instance, the key and (optionally) the IV
	ent, err := entities.NewAES256EncrypterEntity("ID", bccspInst, encKey, IV)
	if err != nil {
		return err
	}

	// here, we encrypt cleartextValue and assign it to key
	return encryptAndPutState(stub, ent, key, cleartextValueAsByte)
}

// Decrypt exposes how to read from the ledger and decrypt using an AES 256
// bit key that has been provided to the chaincode through the transient field.
func Decrypt(stub shim.ChaincodeStubInterface, key string, decKey, IV []byte) ([]byte, error) {
	// create the encrypter entity - we give it an ID, the bccsp instance, the key and (optionally) the IV
	ent, err := entities.NewAES256EncrypterEntity("ID", bccspInst, decKey, IV)
	if err != nil {
		return nil, err
	}

	// here we decrypt the state associated to key
	return getStateAndDecrypt(stub, ent, key)
}
