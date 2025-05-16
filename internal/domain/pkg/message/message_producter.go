package message

import (
	"context"
	shared "go-cs/internal/pkg/domain"
	"go-cs/pkg/bus"

	"github.com/google/uuid"
)

const DomainMessagePublishEventType string = "domain.message.publish_event"

type DomainMessagePublishEvent struct {
	Id       string
	Messages shared.DomainMessages
	Ctx      context.Context
}

type DomainMessageProducer struct {
}

func NewDomainMessageProducer() *DomainMessageProducer {
	return &DomainMessageProducer{}
}

func (d *DomainMessageProducer) Send(ctx context.Context, messages shared.DomainMessages) {
	go func() {
		//打包发送
		bus.Emit(DomainMessagePublishEventType, DomainMessagePublishEvent{
			Id:       uuid.NewString(),
			Messages: messages,
			Ctx:      context.WithoutCancel(ctx),
		})
	}()
}
