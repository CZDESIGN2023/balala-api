package sms

import "testing"

const (
	TEST_ACCESS_KEY_ID     = ""
	TEST_ACCESS_KEY_SECRET = ""
	TEST_SIGN_NAME         = ""
)

func TestClient_Send(t *testing.T) {
	cli := NewClient(TEST_ACCESS_KEY_ID, TEST_ACCESS_KEY_SECRET, TEST_SIGN_NAME)

	var phoneNumber = ""
	var code = "123456"
	var templateCode = ""

	err := cli.Send(phoneNumber, code, templateCode)
	if err != nil {
		t.Error(err)
	}
}
