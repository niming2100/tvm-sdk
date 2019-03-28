package util

import (
	"archive/tar"
	"compress/gzip"
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
	if strings.EqualFold(md5Str, fileMD5) {
		return true, nil;
	} else {
		return false, errors.Errorf("filemd5 check not pass")
	}
}

func compress(file *os.File, prefix string, tw *tar.Writer) error {

	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func Compress(path string, dest string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	var files = []*os.File{file}

	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		err := compress(file, "", tw)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func ModifyPathUserGroup(path string, uid int, gid int) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	var files = []*os.File{file}

	for _, file := range files {
		err := modifyUserGroup(file, uid, gid)
		if err != nil {
			return err
		}
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func modifyUserGroup(file *os.File, uid int, gid int) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	// 改变文件所有者
	err = os.Chown(file.Name(), uid, gid)
	if err != nil {
		return err
	}
	if info.IsDir() {
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = modifyUserGroup(f, uid, gid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
