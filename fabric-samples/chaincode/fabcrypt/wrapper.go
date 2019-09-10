package fabcrypt

import (
	"crypto/md5"
	"encoding/hex"

	// "fmt"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
)

// EncCC example simple Chaincode implementation of a chaincode that uses encryption/signatures
type EncCC struct {
	bccspInst bccsp.BCCSP
}

// Encrypter encrypts the data and writes to the ledger
func Encrypter(stub shim.ChaincodeStubInterface, key string, valueAsBytes []byte) error {
	factory.InitFactories(nil)
	encCC := EncCC{factory.GetDefault()}
	encKey := make([]byte, 32)
	iv := make([]byte, 16)
	ent, err := entities.NewAES256EncrypterEntity("ID", encCC.bccspInst, encKey, iv)
	if err != nil {
		return err
	}

	return encryptAndPutState(stub, ent, key, valueAsBytes)
}

// Decrypter decrypts the data and writes to the ledger
func Decrypter(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	factory.InitFactories(nil)
	encCC := EncCC{factory.GetDefault()}
	// fmt.Println(encCC)
	decKey := make([]byte, 32)
	iv := make([]byte, 16)
	ent, err := entities.NewAES256EncrypterEntity("ID", encCC.bccspInst, decKey, iv)
	if err != nil {
		return nil, err
	}

	return getStateAndDecrypt(stub, ent, key)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
