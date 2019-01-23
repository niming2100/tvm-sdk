package trans

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"os"
	"tvm-sdk/blockchain"
	triasConf "tvm-sdk/config"
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
	var filePath = os.Getenv("GOPATH") + "/src" + triasConf.GetContractPath() + "/" + request.GetUser() + "/" + request.GetAddress() + "/" + request.GetContractName() + "/";
	var fileName = request.GetContractName() + fileSuffix;
	isExists, err := t_utils.PathExists(filePath + fileName);
	if err != nil {
		fmt.Println("checkFilePathFails", err);
		return returnErrorResponse(err);
	}

	if !isExists {
		err := t_utils.FileDownLoad(filePath, fileName, triasConf.GetIPFSAddress()+request.GetAddress());
		if err != nil {
			fmt.Println("Download contract happens a error", err);
			return returnErrorResponse(err);
		}
	}
	if _, err := t_utils.CheckFileMD5(filePath+fileName, request.GetCheckMD5()); err != nil {
		return returnErrorResponse(err);
	}

	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: triasConf.GetOrderServer(),

		// Channel parameters
		ChannelID:     triasConf.GetChannelID(),
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric/singlepeer/channel-artifacts/mychannel.tx",

		// Chaincode parameters
		ChainCodeID:      request.GetContractName(),
		ChainCodeVersion: request.GetContractVersion(),
		ChaincodeGoPath:  os.Getenv("GOPATH"),
		ChaincodePath:    triasConf.GetContractPath() + filePath[len(os.Getenv("GOPATH") + "/src" + triasConf.GetContractPath()):],
		OrgAdmin:         triasConf.GetOrgAdmin(),
		OrgName:          triasConf.GetOrgName(),
		ConfigFile:       "config-sdk.yaml",

		// User parameters
		UserName: triasConf.GetUserName(),
	}

	fSetup.Initialize()

	defer fSetup.CloseSDK()


	funName,args,_ := getFuncAndArgs(request.GetCommand());

	var respStr = ""
	var respErr = error(nil)
	switch request.Operation {
	case "instantiate":
		result,err := fSetup.InstantiateCC(t_utils.StringArrayToByte(args))
		respStr = result
		respErr = err
		break;
	case "install":
		result,err := fSetup.InstallCC()
		respStr = result
		respErr = err
		break;
	case "query":
		result,err := fSetup.Query(funName,t_utils.StringArrayToByte(args[1:]))
		respStr = result
		respErr = err
		break;
	case "invoke":
		result,err := fSetup.Invoke(funName,t_utils.StringArrayToByte(args[1:]))
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


func getFuncAndArgs(command string)(string,[]string,error){
	argMap := make(map[string]interface{});
	err := json.Unmarshal([]byte(command), &argMap)
	if(err!=nil){
		return "",nil,err
	}
	if(argMap["function"]==nil){// func in args
		var data map[string][]string;
		if err := json.Unmarshal([]byte(command), &data); err == nil {
			stringArray := data["Args"];
			var funcName = stringArray[0];
			return funcName,stringArray,nil;
		} else {
			return "",nil,err;
		}
	}else{// func is Independent
		//funcString,_ := argMap["function"].(string);
		//args
	}


	return "",nil,nil
}
