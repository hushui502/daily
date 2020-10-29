package main

import (
	"fabric"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

func TestFunc(t *testing.T) {
	cc := new(fabric.SimpleChaincode)

	stub := shimtest.NewMockStub("sacc", cc)

	stub.MockInit("1", [][]byte{[]byte("a"), []byte("90")})

	res := stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("a")})

	fmt.Println("The value of a is", string(res.Payload))
}
