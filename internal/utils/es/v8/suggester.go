package v8

// Represents the generic suggester interface.
// A suggester's only purpose is to return the
// source of the query as a JSON-serializable
// object. Returning a map[string]interface{}
// will do.
type Suggester interface {
	Name() string
	Source(includeName bool) (interface{}, error)
}
