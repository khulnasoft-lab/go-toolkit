package inject

import (
	"fmt"
	"reflect"
)

// Container is the interface used by the simple injection container
type Container interface {
	// Register registers a provider function
	Register(providers ...any)

	// Bind binds a value directly into the container
	Bind(values ...any)

	// Resolve returns the resolved value for a given type, e.g. value, err := c.Resolve(Type{})
	Resolve(typ any) (any, error)

	// Invoke invokes a function with injected parameters and provided, ordered parameters that may optionally return a
	// value and optionally return a single an error as the last function signature
	Invoke(fn any, args ...any) (any, error)
}

// UnresolvedError is a standard error that can be returned from a container.Get call
type UnresolvedError struct {
	typ reflect.Type
}

func (n UnresolvedError) Error() string {
	return fmt.Sprintf("unresolved: '%s'", typeName(n.typ))
}

var _ error = (*UnresolvedError)(nil)

// Unresolved can be used to check if a Get call returns a "not found" error
var Unresolved = UnresolvedError{}
