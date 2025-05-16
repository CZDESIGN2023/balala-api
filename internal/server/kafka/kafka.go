package kafka

import (
	"context"
	"go-cs/internal/conf"
	"go-cs/internal/pkg/openapi"
	"go-cs/pkg/kratos-transport/broker"
	"go-cs/pkg/kratos-transport/transport/kafka"

	"github.com/go-kratos/kratos/v2/log"
)

type TestMessage struct {
	Channel string `json:"Channel"`
	Message string `json:"Message"`
}

var Server *kafka.Server

func handleKafkaMessage(_ context.Context, topic string, headers broker.Headers, msg *TestMessage) error {
	log.Infof("Topic %s, Headers: %+v, Payload: %+v", topic, headers, msg)
	return nil
}

func NewKafkaBroker(c *conf.Data, open *openapi.OpenAPIService, logger log.Logger) *kafka.Server {
	Server = kafka.NewServer(
		kafka.WithAddress(c.Kafka.Source),
		kafka.WithCodec("json"),
	)
	Server.Init()
	if err := Server.Connect(); err != nil {
		log.Infof("cant connect to broker, skip: %v", err)
	}

	// subscribe
	ctx := context.Background()
	topic := "server_user"
	queue := "chat-group"
	disableAutoAck := false

	Server.RegisterSubscriber(ctx, topic, queue, disableAutoAck,
		func(ctx context.Context, event broker.Event) error {
			log.Infof("Topic %s, Payload: %+v", topic, event.Message().Body)
			open.RouteToOther(ctx, event.Message().Body)
			return nil
		},
		func() broker.Any { return nil },
		broker.WithQueueName(queue),
	)
	return Server
}
