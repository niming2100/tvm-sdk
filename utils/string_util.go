package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	triasConf "tvm-sdk/config"
)

func StringArrayToByte(strArray []string) [][]byte {
	var args [][]byte;
	for _, v := range strArray {
		args = append(args, []byte(v));
	}
	return args;
}

func Sha256(message string) string {
	key := []byte(triasConf.Secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func GetFuncAndArgs(command string) (string, []string, error) {
	argMap := make(map[string]interface{});
	err := json.Unmarshal([]byte(command), &argMap)
	if (err != nil) {
		return "", nil, err
	}
	if (argMap["function"] == nil) { // func in args
		var data map[string][]string;
		if err := json.Unmarshal([]byte(command), &data); err == nil {
			stringArray := data["Args"];
			var funcName = stringArray[0];
			return funcName, stringArray, nil;
		} else {
			return "", nil, err;
		}
	}
	return "", nil, nil
}
