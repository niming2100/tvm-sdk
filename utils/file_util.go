package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FileDownLoad(localPath string, fileName string, urlPath string) error {
	fmt.Println("start download code", urlPath)
	res, err := http.Get(urlPath);
	if err != nil {
		return err;
	}
	os.MkdirAll(localPath, os.ModePerm)
	f, err := os.Create(localPath + fileName)

	if err != nil {
		return err;
	}
	io.Copy(f, res.Body)
	fmt.Println("end download code", urlPath)
	return nil;
}

func CheckFileMD5(filePath string, fileMD5 string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	md5Hash := md5.New();
	if _, err := io.Copy(md5Hash, f); err != nil {
		return false, err
	}
	fmt.Printf("%x\n", md5Hash.Sum(nil))
	md5Str := hex.EncodeToString(md5Hash.Sum(nil))
	if strings.EqualFold(md5Str,fileMD5){
		return true, nil;
	}else {
		return false,errors.Errorf("filemd5 check not pass")
	}
}
