package main

import (
	"crypto/rand"
	"fmt"
	"go-cs/internal/utils"
)

// 產生密碼的小工具
// 可以把產出的密碼再手動替換幾個特殊字符增加強度
func main() {
	for i := 0; i < 5; i++ {
		key, err := generateKey(32) // 32 characters
		if err != nil {
			fmt.Println("Error generating key:", err)
			return
		}
		fmt.Println(key)
	}
}

func generateKey(length int) (string, error) {
	bytes := make([]byte, length*2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	encoded := utils.Base62Encode(bytes)
	var result []rune
	for i, r := range encoded {
		if i%2 != 0 {
			result = append(result, r)
		}
		if len(result) == length {
			break
		}
	}
	return string(result), nil
}
