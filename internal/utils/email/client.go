package email

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"time"
)

type Client struct {
	host     string
	account  string
	password string
	pool     *email.Pool
}

func NewClient(host, account, password string) (*Client, error) {
	pool, err := email.NewPool(
		host,
		4,
		smtp.PlainAuth("", account, password, host),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		host:     host,
		account:  account,
		password: password,
		pool:     pool,
	}, nil
}

func (c *Client) Send(subject, code string, to []string) error {
	e := &email.Email{
		To:      to,
		From:    c.account,
		Subject: subject,
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte(fmt.Sprintf(HTML_TEMPLATE, code)),
		Headers: textproto.MIMEHeader{},
	}

	err := c.pool.Send(e, 10*time.Second)
	if err != nil {
		return err
	}
	return nil
}
