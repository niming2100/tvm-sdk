package validate

import (
	"tvm-sdk/proto/tm"
)

func RequestValidate(request *tm.ExecuteContractRequest) (isCorect bool, err error) {
	return true, nil;
}
