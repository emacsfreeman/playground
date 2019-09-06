// DISCLAIMER:
// THIS SAMPLE CODE MAY BE USED SOLELY AS PART OF THE TEST AND EVALUATION OF THE ALMERYS/BE|YS
// BLOCKCHAIN SERVICE (THE "SERVICE") AND IN ACCORDANCE WITH THE TERMS OF THE AGREEMENT FOR THE SERVICE.
// THIS SAMPLE CODE PROVIDED "AS IS", WITHOUT ANY WARRANTY, ESCROW, TRAINING, MAINTENANCE, OR SERVICE
// OBLIGATIONS WHATSOEVER ON THE PART OF ALMERYS/BE|YS.

package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/attrmgr"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

type Chaincode struct {
}

type ChaincodeStubInterface interface {
	// GetCreator returns `SignatureHeader.Creator` (e.g. an identity)
	// of the `SignedProposal`. This is the identity of the agent (or user)
	// submitting the transaction.
	GetCreator() ([]byte, error)
}

type ClientIdentity interface {

	// GetID returns the ID associated with the invoking identity.  This ID
	// is guaranteed to be unique within the MSP.
	GetID() (string, error)

	// Return the MSP ID of the client
	GetMSPID() (string, error)

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

// ClientIdentityImpl implements the ClientIdentity interface
type clientIdentityImpl struct {
	stub  ChaincodeStubInterface
	mspID string
	cert  *x509.Certificate
	attrs *attrmgr.Attributes
}

// Unmarshals the bytes returned by ChaincodeStubInterface.GetCreator method and
// returns the resulting msp.SerializedIdentity object
func (c *clientIdentityImpl) getIdentity() (*msp.SerializedIdentity, error) {
	sid := &msp.SerializedIdentity{}
	creator, err := c.stub.GetCreator()
	if err != nil || creator == nil {
		return nil, errors.WithMessage(err, "failed to get transaction invoker's identity from the chaincode stub")
	}
	err = proto.Unmarshal(creator, sid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal transaction invoker's identity")
	}
	return sid, nil
}

// Initialize the client
func (c *clientIdentityImpl) init() error {
	signingID, err := c.getIdentity()
	if err != nil {
		return err
	}
	c.mspID = signingID.GetMspid()
	idbytes := signingID.GetIdBytes()
	block, _ := pem.Decode(idbytes)
	if block == nil {
		return errors.New("Expecting a PEM-encoded X509 certificate; PEM block not found")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.Wrap(err, "failed to parse certificate")
	}
	c.cert = cert
	attrs, err := attrmgr.New().GetAttributesFromCert(cert)
	if err != nil {
		return errors.WithMessage(err, "failed to get attributes from the transaction invoker's certificate")
	}
	c.attrs = attrs
	return nil
}

// AssertAttributeValue checks to see if an attribute value equals the specified value
func AssertAttributeValue(stub ChaincodeStubInterface, attrName, attrValue string) error {
	c, err := New(stub)
	if err != nil {
		return err
	}
	return c.AssertAttributeValue(attrName, attrValue)
}

// New returns an instance of ClientIdentity
func New(stub ChaincodeStubInterface) (ClientIdentity, error) {
	c := &clientIdentityImpl{stub: stub}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetID returns the ID associated with the invoking identity.  This ID
// is guaranteed to be unique within the MSP.
func GetID(stub ChaincodeStubInterface) (string, error) {
	c, err := New(stub)
	if err != nil {
		return "", err
	}
	return c.GetID()
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
		if stub.GetID == "allowed" {
			return write(stub, args)
		} else {
			return shim.Error("not allowed to write")
		}
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

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(Chaincode)); err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}
