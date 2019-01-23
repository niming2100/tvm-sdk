package contract

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
	"strings"
	tvm_conf "tvm-light/config"
)

type Contract struct {
	peerAddress     string
	contractName    string
	contractType    string
	contractPath    string
	contractVersion string
	channelID       string
	orgName         string
	args            string
	action          string
}

type Args struct {
	Args interface{} `json:"Args"`
}

// {\"Args\":[\"query\",\"a\"]}
const (
	CMD_DOCKER = "docker"
)

var docker_command []string = []string{"exec", "cli"}

func NewContract(peerAddress string, contractName string, contractType string, contractPath string, contractVersion string, channelID string, orgName string, args string, action string) *Contract {
	return &Contract{peerAddress: peerAddress, contractName: contractName, contractType: contractType, contractPath: contractPath, contractVersion: contractVersion, channelID: channelID, orgName: orgName, args: args, action: action}
}


func (c *Contract) RunContract() (string, error) {

	var resp string;
	var err error = nil;
	switch c.action {
	case "instantiate":
		resp,err = c.instantiate()
		break;
	case "install":
		resp,err =  c.InstallContract()
	default:
		resp,err =  c.execute()
		break;
	}
	return resp,err;
}

func (c *Contract) InstallContract() (string,error) {
	var filePath string = tvm_conf.GetDockerPath() + c.contractPath[len(tvm_conf.GetContractPath()):];
	fmt.Println("dockerPath:",tvm_conf.GetDockerPath())
	params := []string{"peer", "chaincode", "install", "-n", c.contractName, "-p", filePath, "-v", c.contractVersion};
	return runCommand(params);
}

func (c *Contract) instantiate() (string, error) {
	params := []string{"peer", "chaincode", "instantiate", "-o", tvm_conf.GetOrderServer(), "-C", c.channelID, "-n", c.contractName, "-v", c.contractVersion, "-c", c.args};
	return runCommand(params);
}

func (c *Contract) execute() (string, error) {
	// peer chaincode $action -C $channelID -n $cname -c $args
	params := []string{"peer", "chaincode", c.action, "-C", c.channelID, "-n", c.contractName, "-c", c.args};
	return runCommand(params);

}

func runCommand(param []string)(string, error){
	cmd := exec.Command(CMD_DOCKER, append(docker_command, param...)...);
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	stderrOut, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 保证关闭输出流
	defer stdout.Close()
	defer stderrOut.Close()
	// 运行命令
	fmt.Println(cmd.Args)
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return "", err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	errBytes, err := ioutil.ReadAll(stderrOut)

	var opString string = string(opBytes)
	var errString string = string(errBytes)
	cmd.Wait()
	if !strings.EqualFold(opString,"") {
		fmt.Println(opString)
		return opString,nil
	} else {
		fmt.Println(errString)
		err = solveErrorResult(errString)
		if(err != nil){
			return "install error",err
		} else {
			return "",err
		}
	}
}

func solveErrorResult(result string) error{
	var err error = nil;
	if(strings.Contains(result,"Error: ")){
		err = errors.Errorf("Execute contract happens a error!")
	}
	return err
}