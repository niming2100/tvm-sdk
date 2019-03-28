package trans

import (
	"fmt"
	"golang.org/x/net/context"
	"os"
	"tvm-sdk/blockchain"
	triasConf "tvm-sdk/config"
	"tvm-sdk/contract"
	"tvm-sdk/proto/tm"
	t_utils "tvm-sdk/utils"
	"tvm-sdk/validate"
)

func NewMWService() *server {
	return &server{}
}

const (
	fileSuffix = ".go"
)

type server struct {
}

func (serv *server) ExecuteContract(ctx context.Context, request *tm.ExecuteContractRequest) (*tm.ExecuteContractResponse) {
	// TODO validate
	isCorect, err := validate.RequestValidate(request);
	if !isCorect || err != nil {
		fmt.Println("Contract validate fails", err);
		return returnErrorResponse(err);
	}
	// TODO CheckContract is install
	var filePath = triasConf.TriasConfig.ContractPath + "/" + request.GetUser() + "/" + request.GetAddress() + "/" + request.GetContractName() + "/";
	var fileName = request.GetContractName() + fileSuffix;
	isExists, err := t_utils.PathExists(filePath + fileName);
	if err != nil {
		fmt.Println("checkFilePathFails", err);
		return returnErrorResponse(err);
	}

	if !isExists {
		err := t_utils.FileDownLoad(filePath, fileName, triasConf.TriasConfig.IPFSAddress+request.GetAddress());
		if err != nil {
			fmt.Println("Download contract happens a error", err);
			return returnErrorResponse(err);
		}
	}
	if _, err := t_utils.CheckFileMD5(filePath+fileName, request.GetCheckMD5()); err != nil {
		return returnErrorResponse(err);
	}

	fmt.Println(filePath)

	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: triasConf.TriasConfig.OrderServer,

		// Channel parameters
		ChannelID:     triasConf.TriasConfig.ChannelID,
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric/singlepeer/channel-artifacts/mychannel.tx",

		// Chaincode parameters
		ChainCodeID:      request.GetContractName(),
		ChainCodeVersion: request.GetContractVersion(),
		ChaincodeGoPath:  os.Getenv("GOPATH"),
		ChaincodePath:    triasConf.TriasConfig.DockerPath + filePath[len(triasConf.TriasConfig.ContractPath):],
		OrgAdmin:         triasConf.TriasConfig.OrgAdmin,
		OrgName:          triasConf.TriasConfig.OrgName,
		ConfigFile:       os.Getenv("GOPATH") + "/src/tvm-sdk/config_e2e_single_org.yaml",

		// User parameters
		UserName: triasConf.TriasConfig.OrgAdmin,
	}

	fSetup.Initialize()

	defer fSetup.CloseSDK()

	funName, args, _ := t_utils.GetFuncAndArgs(request.GetCommand());
	var respStr = ""
	var respErr = error(nil)
	switch request.Operation {
	case "instantiate":
		result, err := fSetup.InstantiateCC(t_utils.StringArrayToByte(args))
		respStr = result
		respErr = err
		break;
	case "install":
		result, err := fSetup.InstallCC()
		respStr = result
		respErr = err
		break;
	case "query":
		result, err := fSetup.Query(funName, t_utils.StringArrayToByte(args[1:]))
		respStr = result
		respErr = err
		break;
	case "invoke":
		result, err := fSetup.Invoke(funName, t_utils.StringArrayToByte(args[1:]))
		respStr = result
		respErr = err
		break;
	default:
		//error
		break;
	}

	if respErr != nil {
		fmt.Println(respErr)
		return returnErrorResponse(respErr);
	} else {
		c_hash := calculateHash(request)
		contract.UpdateCurrentHash(triasConf.BasicHashKey, c_hash)
	}

	resp := &tm.ExecuteContractResponse{
		Code:    1,
		Data:    respStr,
		Message: "success",
	}

	return resp;
}

func returnErrorResponse(err error) (*tm.ExecuteContractResponse) {
	resp := &tm.ExecuteContractResponse{
		Code:    -1,
		Data:    err.Error(),
		Message: "fail",
	}
	return resp;
}

func calculateHash(request *tm.ExecuteContractRequest) string {
	var message string = string(request.GetContractName() + request.GetUser() + request.GetOperation() + string(request.GetTimestamp()) + request.GetCheckMD5());
	return message;
}
