package utils

import (
	"math/rand"
)

func GenerateVerifyCode(n int) string {
	var s = make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = byte(rand.Intn(10) + '0')
	}

	return string(s)
}
