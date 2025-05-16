package rand

import (
	"crypto/rand"
	"encoding/binary"
)

// RandomInt64 随机64位整数
func RandomInt64() int64 {
	var randNum int64
	err := binary.Read(rand.Reader, binary.BigEndian, &randNum)
	if err != nil {
		return 0
	}
	return randNum
}

// RandomInt32 随机32位
func RandomInt32() int32 {
	var randNum int32
	err := binary.Read(rand.Reader, binary.BigEndian, &randNum)
	if err != nil {
		return 0
	}
	return randNum
}
