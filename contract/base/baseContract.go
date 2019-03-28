package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type BasicTrias struct {
}

func (b *BasicTrias) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (b *BasicTrias) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "invoke" {
		return b.invoke(stub, args)
	} else if fn == "query" {
		return b.query(stub, args)
	}
	return shim.Error("Invoke function error")
}

func (b *BasicTrias) query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var key = args[0]
	dataBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(dataBytes)
}

func (b *BasicTrias) invoke(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	var key, hash string
	key = args[0]
	hash = args[1]
	if hash == "" {
		return shim.Error("CurrentHash can not be empty")
	}
	err := stub.PutState(key, []byte(hash))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

func main() {
	err := shim.Start(new(BasicTrias))
	if err != nil {
		fmt.Println("start error")
	}
}
