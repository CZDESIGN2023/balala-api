package errs

import (
	"go-cs/api/comm"
	"testing"
)

func Test(t *testing.T) {
	err := Business(nil, comm.ErrorCode_ERROR_MSG_TEST)
	t.Log(err)
}
