package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	api "go-cs/pkg/kratos-transport/_example/api/manual"
	"go-cs/pkg/kratos-transport/broker"
	rabbitmqBroker "go-cs/pkg/kratos-transport/broker/rabbitmq"
	"go-cs/pkg/kratos-transport/transport/rabbitmq"
)

const (
	testBroker = "amqp://user:bitnami@127.0.0.1:5672"

	testExchange = "test_exchange"
	testQueue    = "test_queue"
	testRouting  = "test_routing_key"
)

func handleHygrothermograph(_ context.Context, topic string, headers broker.Headers, msg *api.Hygrothermograph) error {
	log.Infof("Topic %s, Headers: %+v, Payload: %+v\n", topic, headers, msg)
	return nil
}

func main() {
	ctx := context.Background()

	rabbitmqSrv := rabbitmq.NewServer(
		rabbitmq.WithAddress([]string{testBroker}),
		rabbitmq.WithCodec("json"),
		rabbitmq.WithExchange(testExchange, true),
	)

	_ = rabbitmqSrv.RegisterSubscriber(ctx,
		testRouting,
		api.RegisterHygrothermographJsonHandler(handleHygrothermograph),
		api.HygrothermographCreator,
		broker.WithQueueName(testQueue),
		rabbitmqBroker.WithDurableQueue())

	app := kratos.New(
		kratos.Name("rabbitmq"),
		kratos.Server(
			rabbitmqSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
