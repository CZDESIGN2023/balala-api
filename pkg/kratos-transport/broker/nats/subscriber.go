package nats

import (
	natsGo "github.com/nats-io/nats.go"
	"go-cs/pkg/kratos-transport/broker"
)

type subscriber struct {
	s    *natsGo.Subscription
	opts broker.SubscribeOptions
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.s.Subject
}

func (s *subscriber) Unsubscribe() error {
	return s.s.Unsubscribe()
}
