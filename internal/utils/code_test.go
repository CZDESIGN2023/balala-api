package utils

import "testing"

func TestGenerateVerifyCode(t *testing.T) {
	code := GenerateVerifyCode(6)

	t.Log(code)
}
