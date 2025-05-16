package message

import (
	"go-cs/internal/utils"
	"go-cs/pkg/bus"
	"reflect"
	"runtime/debug"

	"github.com/go-kratos/kratos/v2/log"
)

type DomainMessageConsumer struct {
	log *log.Helper
	// listeners map[string]func(DomainMessagePublishEvent)
}

func NewDomainMessageConsumer(
	logger log.Logger,
) *DomainMessageConsumer {

	moduleName := "domain_message_consumer"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &DomainMessageConsumer{
		// listeners: make(map[string]func(DomainMessagePublishEvent), 0),
		log: hlog,
	}
}

func (d *DomainMessageConsumer) SetMessageListener(id string, listener func(DomainMessagePublishEvent)) {

	listenerId := "domain_message_consumer." + id
	bus.On(DomainMessagePublishEventType, listenerId, func(args ...any) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					d.log.Errorf("handle event %s panic: %v, %s", DomainMessagePublishEventType, err, debug.Stack())
				}
			}()

			var rvs []reflect.Value
			for _, v := range args {
				rvs = append(rvs, reflect.ValueOf(v))
			}
			reflect.ValueOf(listener).Call(rvs)
		}()
	})
}
