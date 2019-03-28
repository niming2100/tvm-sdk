package util

import (
	"github.com/ipfs/go-ipfs-api"
	"io"
	"os"
	t_conf "tvm-sdk/config"
)

func AddFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	var sh *shell.Shell
	sh = shell.NewShell(t_conf.TriasConfig.IpfsAPIAddress)
	hash, err := sh.Add(f)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func GetFile(hash string, dest string) error {
	var sh *shell.Shell
	sh = shell.NewShell(t_conf.TriasConfig.IpfsAPIAddress)
	read, err := sh.Cat(hash)
	if err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	io.Copy(f, read)
	return nil;
}
