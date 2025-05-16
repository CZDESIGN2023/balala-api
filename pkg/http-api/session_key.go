package http_api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const (
	deaufulSaltKey = "97875k28c0d73109"
)

// MakeSign 生成
func MakeSign(info map[string]interface{}, saltKey string) (string, error) {
	if info == nil {
		return "", fmt.Errorf("make session fail, because of no info")
	}
	if saltKey == "" {
		saltKey = deaufulSaltKey
	}

	info["tm"] = time.Now().Unix()
	str, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	qsX := base64.StdEncoding.EncodeToString([]byte(str))

	hash := md5.Sum([]byte(qsX + saltKey))
	sign := hex.EncodeToString(hash[:])

	return "info=" + qsX + "&sign=" + sign, nil
}

// CheckSign 解开
func CheckSign(info string, sign string, saltKey string) (map[string]interface{}, error) {
	if info == "" {
		return nil, fmt.Errorf("check session fail, because of no info")
	}
	if saltKey == "" {
		saltKey = deaufulSaltKey
	}

	hash := md5.Sum([]byte(info + saltKey))
	md5Str := hex.EncodeToString(hash[:])
	if md5Str != sign {
		return nil, fmt.Errorf("sessionKey check error sign")
	}

	appInfoStr, err := base64.StdEncoding.DecodeString(info)
	if err != nil {
		return nil, fmt.Errorf("sessionKey check base64 error")
	}

	var appInfo map[string]interface{}
	if err := json.Unmarshal(appInfoStr, &appInfo); err != nil {
		return nil, err
	}

	return appInfo, nil
}
