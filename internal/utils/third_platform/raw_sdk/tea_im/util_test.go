package tea_im

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cast"
	"go-cs/pkg/stream"
	"testing"
)

func Test_sign(t *testing.T) {

	params := map[string]any{
		"chat_token":  "YmFsYWxhICAgICAgICAyMjE3MTkzMTI3NDM2NTU2Mw==",
		"sign":        "8a182cdc1fa23471694616fdd9511c41",
		"id":          1,
		"pf_code":     "balala",
		"user_name":   "一需要吕不韦",
		"user_id":     22,
		"create_time": 1719312743,
	}

	s := params["sign"]
	delete(params, "sign")

	paramsM := stream.MapValue(params, func(v any) string {
		return cast.ToString(v)
	})

	priKey := "iTMifz7PDBkBCre4B3meNhw5z7tBAAKi"

	arg := makeArgs(paramsM) + fmt.Sprintf("&pri_key=%v", priKey)

	t.Log(arg)

	md5Hash := md5.Sum([]byte(arg))
	s2 := hex.EncodeToString(md5Hash[:])

	t.Log(s == s2)

	t.Log(s)
	t.Log(s2)
}
