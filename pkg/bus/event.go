package bus

import (
	"go-cs/pkg/sync1"
	"reflect"
)

type HandFunc any

type ID any

type event struct {
	name any
	Ids  sync1.Map[ID, HandFunc]
}

type Manager struct {
	events sync1.Map[any, *event]
}

func (m *Manager) On(eventName any, identity ID, handFunc HandFunc) {
	of := reflect.TypeOf(handFunc)
	if of.Kind() != reflect.Func {
		panic("must be func")
	}

	t, _ := m.events.LoadOrStore(eventName, &event{name: eventName})
	t.Ids.Store(identity, handFunc)
}

func (m *Manager) Off(eventName any, identity ID) {
	t, loaded := m.events.Load(eventName)
	if !loaded {
		return
	}

	t.Ids.Delete(identity)
}

func (m *Manager) Emit(eventName any, args ...any) {
	t, loaded := m.events.Load(eventName)
	if !loaded {
		return
	}

	var rvs []reflect.Value
	for _, v := range args {
		rvs = append(rvs, reflect.ValueOf(v))
	}

	t.Ids.Range(func(k ID, fn HandFunc) {
		reflect.ValueOf(fn).Call(rvs)
	})

	return
}

func (m *Manager) EmitTo(eventName any, identity ID, args ...any) {
	t, loaded := m.events.Load(eventName)
	if !loaded {
		return
	}

	fn, loaded := t.Ids.Load(identity)
	if !loaded {
		return
	}

	var rvs []reflect.Value
	for _, v := range args {
		rvs = append(rvs, reflect.ValueOf(v))
	}
	reflect.ValueOf(fn).Call(rvs)

	return
}
