package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (setup *FabricSetup) Invoke(funcName string, args [][]byte) (string, error) {

	// Prepare arguments

	eventID := "eventInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in trias invoke")

	reg, _, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{
		ChaincodeID:  setup.ChainCodeID,
		Fcn:          funcName,
		Args:         args,
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	return string(response.TransactionID), nil
}
