package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (setup *FabricSetup) Query(funcName string, args [][]byte) (string, error) {

	response, err := setup.client.Query(channel.Request{
		ChaincodeID: setup.ChainCodeID,
		Fcn:         funcName,
		Args:        args,
	}, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}
