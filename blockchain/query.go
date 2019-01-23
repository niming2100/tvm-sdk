package blockchain

import (
	"encoding/json"
	"fmt"

	"citizens/common"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// 查询用户
func (setup *FabricSetup) Query(funcName string,args [][]byte) (string, error) {



	response, err := setup.client.Query(channel.Request{
		ChaincodeID: setup.ChainCodeID,
		Fcn:         funcName,
		Args:        args,
	})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	people := common.People{}
	err = json.Unmarshal([]byte(response.Payload), &people)
	fmt.Println(people)

	return string(response.Payload), nil
}
