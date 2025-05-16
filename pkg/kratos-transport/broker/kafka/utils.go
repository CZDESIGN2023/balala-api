package kafka

import (
	kafkaGo "github.com/segmentio/kafka-go"

	"go-cs/pkg/kratos-transport/broker"
)

func kafkaHeaderToMap(h []kafkaGo.Header) broker.Headers {
	m := broker.Headers{}
	for _, v := range h {
		m[v.Key] = string(v.Value)
	}
	return m
}
