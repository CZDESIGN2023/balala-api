package bus

var global = New()

func New() *Manager {
	return &Manager{}
}

func On(event any, id ID, handFunc HandFunc) {
	global.On(event, id, handFunc)
}

func Off(event any, id ID) {
	global.Off(event, id)
}

func Emit(event any, args ...any) {
	global.Emit(event, args...)
}

func EmitTo(event any, id ID, args ...any) {
	global.EmitTo(event, id, args...)
}
