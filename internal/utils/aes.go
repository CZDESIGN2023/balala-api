package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

var aes256DefaultKey []byte

func init() {
	aes256DefaultKey = []byte("c_qM,754/Z@6c.&-23~)98!d:]{|`*f%")
}

func EncryptAES(plaintext, key []byte) (string, error) {
	if key == nil {
		key = aes256DefaultKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	//return url.QueryEscape(base64.URLEncoding.EncodeToString(ciphertext)), nil
	return Base62Encode(ciphertext), nil
}

func DecryptAES(encryptedAES string, key []byte) ([]byte, error) {
	if key == nil {
		key = aes256DefaultKey
	}

	ciphertextBytes, err := Base62Decode(encryptedAES)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return ciphertextBytes, nil
}
