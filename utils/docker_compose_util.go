package util

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	t_conf "tvm-sdk/config"
)

var composeCommoand string = string("docker-compose")

func StartTVM() error {
	command := []string{"-f", t_conf.TriasConfig.GOPATH + t_conf.TriasConfig.ComposeFilePath, "up", "-d"}
	return execCommand(command)
}

func StopTVM() error {
	command := []string{"-f", t_conf.TriasConfig.GOPATH + t_conf.TriasConfig.ComposeFilePath, "stop"}
	return execCommand(command)
}

func RestartTVM() {
	StopTVM()
	StartTVM()
}

func execCommand(command []string) error {
	cmd := exec.Command(composeCommoand, command...);
	stdout, err := cmd.StdoutPipe()
	stderrOut, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	// 保证关闭输出流
	defer stdout.Close()
	defer stderrOut.Close()

	if err := cmd.Start(); err != nil {
		return err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	errBytes, err := ioutil.ReadAll(stderrOut)

	opString := string(opBytes)
	errString := string(errBytes)
	cmd.Wait()
	fmt.Println(opString)
	fmt.Println(errString)
	return nil
}
