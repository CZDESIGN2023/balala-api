package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

var signSecretKey = []byte("J_SvI8IV4$>FRQ7drV.r<V[YWHbgM9gm")

var signSalt = "_UNfj8DyE7epzkLbUPGKfDIbEHkO8YdjY_"

// GenFileSignUri 產生檔案下載的uri, 附加簽名、過期時間(10分鐘)
func GenFileSignUri(path string) string {
	if path == "" {
		return ""
	}

	expiration := time.Now().Add(time.Minute * 10).Unix()
	message := path + strconv.FormatInt(expiration, 10) + signSalt
	sign := ComputeHmac256(message, signSecretKey)

	rawQuery := url.Values{
		"sign": {sign},
		"exp":  {strconv.FormatInt(expiration, 10)},
	}.Encode()

	u := &url.URL{
		Path:     path,
		RawQuery: rawQuery,
	}

	return u.String()
}

// VerifyFileSignUri 驗證前端傳入的uri 是否過期、簽名是否正確
func VerifyFileSignUri(path string, exp string, sign string) (bool, error) {
	// Check if signature and expiration are present
	if path == "" || sign == "" || exp == "" {
		return false, fmt.Errorf("參數錯誤")
	}

	// Check if URL has expired
	expiration, err := strconv.ParseInt(exp, 10, 64)
	if err != nil {
		return false, errors.New("下載鏈結過期")
	}

	if time.Now().Unix() > expiration {
		return false, errors.New("下載鏈結過期")
	}

	// Verify signature
	message := path + exp + signSalt
	expectedSignature := ComputeHmac256(message, signSecretKey)

	if !hmac.Equal([]byte(sign), []byte(expectedSignature)) {
		return false, errors.New("簽名錯誤")
	}

	return true, nil
}

func ComputeHmac256(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
