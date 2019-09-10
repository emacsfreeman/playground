package fabcrypt

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
	"github.com/pkg/errors"
)

// getStateAndDecrypt retrieves the value associated to key,
// decrypts it with the supplied entity and returns the result
// of the decryption
func getStateAndDecrypt(stub shim.ChaincodeStubInterface, ent entities.Encrypter, key string) ([]byte, error) {
	// at first we retrieve the ciphertext from the ledger
	ciphertext, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) == 0 {
		return nil, errors.New("no ciphertext to decrypt")
	}

	return ent.Decrypt(ciphertext)
}

// encryptAndPutState encrypts the supplied value using the
// supplied entity and puts it to the ledger associated to
// the supplied KVS key
func encryptAndPutState(stub shim.ChaincodeStubInterface, ent entities.Encrypter, key string, value []byte) error {
	// at first we use the supplied entity to encrypt the value
	ciphertext, err := ent.Encrypt(value)
	if err != nil {
		return err
	}

	return stub.PutState(key, ciphertext)
}
