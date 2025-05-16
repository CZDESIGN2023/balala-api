package v8

type Query interface {
	Source() (interface{}, error)
}
