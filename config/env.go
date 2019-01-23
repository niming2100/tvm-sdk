package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type TVMconf struct {
	OrderServer    string `yaml:"orderServer"`
	ContractPath   string `yaml:"contractPath"`
	ChannelID      string `yaml:"channelID"`
	OrdererOrgName string `yaml:"ordererOrgName"`
	IPFSAddress    string `yaml:"IPFSAddress"`
	DockerPath     string `yaml:"dockerPath"`
	Port           string `yml:"port"`
	OrgAdmin       string `yml:"orgAdmin"`
	OrgName        string `yml:"orgName"`
	UserName       string `yml:"userName"`
	PeerServer       string `yml:"peerServer"`

}

var triasConfig = TVMconf{}

func init() {
	//var filePath = "/home/Polarbear/workGo/src/tvm-sdk/config.yml"
	var filePath = "config.yml"
	data, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(data, &triasConfig)
}

func GetOrderServer() string {
	return triasConfig.OrderServer;
}
func GetContractPath() string {
	return triasConfig.ContractPath;
}

func GetChannelID() string {
	return triasConfig.ChannelID;
}

func GetOrdererOrgName() string {
	return triasConfig.OrdererOrgName
}

func GetIPFSAddress() string {
	return triasConfig.IPFSAddress
}

func GetDockerPath() string {
	return triasConfig.DockerPath
}

func GetPort() string {
	return triasConfig.Port
}

func GetOrgAdmin() string {
	return triasConfig.OrgAdmin
}

func GetOrgName() string {
	return triasConfig.OrgName
}

func GetUserName() string {
	return triasConfig.UserName
}

func GetPeerServer() string {
	return triasConfig.PeerServer
}