package email

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"testing"
)

const (
	TEST_HOST     = "smtp.qq.com:465"
	TEST_ACCOUNT  = ""
	TEST_PASSWORD = ""
)

func TestClient_Send(t *testing.T) {
	cli, err := NewClient(TEST_HOST, TEST_ACCOUNT, TEST_PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = cli.Send("subject", "code", []string{"example@example.com"})
	if err != nil {
		t.Error(err)
	}
}

func TestSend(t *testing.T) {
	e := &email.Email{
		To:      []string{"test@example.com"},
		From:    "Jordan Wright <test@gmail.com>",
		Subject: "Awesome Subject",
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}

	err := e.Send(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", TEST_ACCOUNT, TEST_PASSWORD, TEST_HOST),
	)
	if err != nil {
		t.Error(err)
	}
}
