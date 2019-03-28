package contract

import (
	"fmt"
	"os"
	"tvm-sdk/blockchain"
	triasConf "tvm-sdk/config"
	t_utils "tvm-sdk/utils"
)

func UpdateCurrentHash(key string, hash string) {
	bSetup := initSdk();
	defer bSetup.CloseSDK()
	c_hash := getValue(key, bSetup);
	final_hash := t_utils.Sha256(hash + triasConf.Secret + c_hash);
	setValue(key, final_hash, bSetup);
}

func UpdateCurrentIpfsAddress(key string, hash string) {
	bSetup := initSdk();
	defer bSetup.CloseSDK()
	setValue(key, hash, bSetup);
}

func GetCurrentHash(key string) string {
	bSetup := initSdk();
	defer bSetup.CloseSDK()
	c_hash := getValue(key, bSetup);
	return c_hash
}

func setValue(key string, value string, setup *blockchain.FabricSetup) {
	var command string = "{\"Args\":[\"invoke\",\"" + key + "\",\"" + value + "\"]}"
	funName, args, _ := t_utils.GetFuncAndArgs(command);
	if _, err := setup.Invoke(funName, t_utils.StringArrayToByte(args[1:])); err != nil {
		fmt.Println(err)
	}
}

func getValue(key string, setup *blockchain.FabricSetup) string {
	var command string = "{\"Args\":[\"query\",\"" + key + "\"]}"
	funName, args, _ := t_utils.GetFuncAndArgs(command);
	if result, err := setup.Query(funName, t_utils.StringArrayToByte(args[1:])); err == nil {
		return result
	} else {
		fmt.Println(err)
		return ""
	}
}

func initSdk() *blockchain.FabricSetup {
	basicSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: triasConf.TriasConfig.OrderServer,

		// Channel parameters
		ChannelID:     triasConf.TriasConfig.ChannelID,
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric/singlepeer/channel-artifacts/mychannel.tx",

		// Chaincode parameters
		ChainCodeID:      triasConf.BasicContractName,
		ChainCodeVersion: triasConf.BasicContractVersion,
		ChaincodeGoPath:  os.Getenv("GOPATH"),
		ChaincodePath:    triasConf.BasicContractPath,
		OrgAdmin:         triasConf.TriasConfig.OrgAdmin,
		OrgName:          triasConf.TriasConfig.OrgName,
		ConfigFile:       os.Getenv("GOPATH") + "/src/tvm-sdk/config_e2e_single_org.yaml",

		// User parameters
		UserName: triasConf.TriasConfig.OrgAdmin,
	}
	basicSetup.Initialize()
	return &basicSetup;
}
