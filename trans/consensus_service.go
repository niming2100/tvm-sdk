package trans

import (
	"encoding/json"
	"golang.org/x/net/context"
	"os"
	triasConf "tvm-sdk/config"
	"tvm-sdk/contract"
	t_utils "tvm-sdk/utils"
)

func NewConsensusService() *consensusServer {
	return &consensusServer{}
}

type consensusServer struct{}

type AsyncTVMRequest struct {
	IpfsHash string `json:"ipfsHash,omitempty"`
}
type uploadResponseBody struct {
	IpfsHash string `json:"ipfsHash"`
}

const (
	tar_suffix = ".tar.gz"
)

func (s *consensusServer) UploadData() *triasConf.CommonResponse {
	// package data
	var filePath string = triasConf.TriasConfig.PackagePath + "/data.tar.gz"
	err := t_utils.Compress(triasConf.TriasConfig.DataPath, filePath)
	if err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// uploadIpfs
	hash, err := t_utils.AddFile(filePath)
	if err != nil {
		return createErrorCommonResponse(err, -1)
	}
	data := uploadResponseBody{
		IpfsHash: hash,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return createErrorCommonResponse(err, -1)
	} else {
		contract.UpdateCurrentIpfsAddress(triasConf.BasicIpfsKey, hash)
	}
	return createSuccessCommonResponse(string(jsonData))
}

func (s *consensusServer) AsyncTVM(ctx context.Context, request *AsyncTVMRequest) *triasConf.CommonResponse {
	// download package from ipfs
	var fileName string = string(triasConf.TriasConfig.PackagePath + request.IpfsHash + tar_suffix)
	if err := t_utils.GetFile(request.IpfsHash, fileName); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// stop docker-compose
	if err := t_utils.StopTVM(); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// clean data files
	if err := os.RemoveAll(triasConf.TriasConfig.DataPath); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// decomporess data file
	if err := t_utils.DeCompress(fileName, triasConf.TriasConfig.DataPath[:len(triasConf.TriasConfig.DataPath)-5]); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// chown 5984.5984
	if err := t_utils.ModifyPathUserGroup(triasConf.TriasConfig.CouchdbInfo.Path, triasConf.TriasConfig.CouchdbInfo.Port, triasConf.TriasConfig.CouchdbInfo.Port); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	// start docker-compose
	if err := t_utils.StartTVM(); err != nil {
		return createErrorCommonResponse(err, -1)
	}
	return createSuccessCommonResponse("")
}

func (s *consensusServer) GetCurrentHash() *triasConf.CommonResponse {
	value := contract.GetCurrentHash(triasConf.BasicHashKey)
	return createSuccessCommonResponse(value)
}

func (s *consensusServer) GetCurrentDataAddress() *triasConf.CommonResponse {
	value := contract.GetCurrentHash(triasConf.BasicIpfsKey)
	return createSuccessCommonResponse(value)
}

func createErrorCommonResponse(err error, code int32) *triasConf.CommonResponse {
	resp := &triasConf.CommonResponse{
		Code:    code,
		Data:    err.Error(),
		Message: "fail",
	}
	return resp;
}

func createSuccessCommonResponse(data string) *triasConf.CommonResponse {
	resp := &triasConf.CommonResponse{
		Code:    1,
		Data:    data,
		Message: "success",
	}
	return resp;
}
