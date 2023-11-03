package inject

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

func NewContainer(parents ...Container) Container {
	c := &container{
		providers: map[reflect.Type]reflect.Value{},
	}
	c.parents = parents
	return c
}

type container struct {
	parents   []Container
	providers map[reflect.Type]reflect.Value
}

var _ Container = (*container)(nil)

func (c *container) Register(providers ...any) {
	for _, provider := range providers {
		v := reflect.ValueOf(provider)
		t := v.Type()
		validateProvider(t, provider)
		// this doesn't handle multiple return values
		rt := t.Out(0)
		bt := baseType(rt)
		c.providers[bt] = v
	}
}

func (c *container) Bind(values ...any) {
	for _, value := range values {
		v := reflect.ValueOf(value)
		t := v.Type()
		bt := baseType(t)
		ft := reflect.FuncOf([]reflect.Type{}, []reflect.Type{t}, false)
		f := reflect.MakeFunc(ft, func(args []reflect.Value) (results []reflect.Value) {
			return []reflect.Value{v}
		})
		c.providers[bt] = f
	}
}

func (c *container) Resolve(typ any) (any, error) {
	var t reflect.Type
	if rt, ok := typ.(reflect.Type); ok {
		t = rt
	} else {
		t = reflect.TypeOf(typ)
	}
	bt := baseType(t)
	if bt == containerInterface {
		return alignReturn(t, reflect.ValueOf(c)), nil
	}
	if p, ok := c.providers[bt]; ok {
		v, err := c.invoke(p)
		if err != nil {
			return nil, err
		}
		return alignReturn(t, v), nil
	}
	for _, p := range c.parents {
		v, err := p.Resolve(t)
		if errors.Is(err, Unresolved) {
			continue
		}
		return v, err
	}
	return nil, UnresolvedError{
		typ: t,
	}
}

func alignReturn(t reflect.Type, v reflect.Value) any {
	vt := v.Type()
	switch {
	case vt.ConvertibleTo(t):
	case isPtr(t) && !isPtr(vt) && !v.CanAddr():
		pv := reflect.New(vt)
		pv.Elem().Set(v)
		v = pv
	case isPtr(t) && !isPtr(vt):
		v = v.Addr()
	case !isPtr(t) && isPtr(vt):
		v = v.Elem()
	}
	return v.Interface()
}

func (c *container) Invoke(fn any, args ...any) (any, error) {
	v := reflect.ValueOf(fn)
	rv, err := c.invoke(v, args...)
	var out any
	if rv.IsValid() && rv.CanInterface() {
		out = rv.Interface()
	}
	return out, err
}

var (
	nilValue           = reflect.ValueOf(nil)
	errorInterface     = reflect.TypeOf((*error)(nil)).Elem()
	containerInterface = reflect.TypeOf((*Container)(nil)).Elem()
)

func (c *container) invoke(fn reflect.Value, args ...any) (reflect.Value, error) {
	t := fn.Type()
	if t.Kind() != reflect.Func {
		return nilValue, fmt.Errorf("unable to invoke non-function type: %s : %+v", t.Name(), fn)
	}

	inArg := 0

	in := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		ft := t.In(i)

		// check positional parameters
		if inArg < len(args) {
			arg := args[inArg]
			argV := reflect.ValueOf(arg)
			if argV.CanConvert(ft) {
				in[i] = argV.Convert(ft)
				inArg++
				continue
			}
		}

		// use type resolution
		v, err := c.Resolve(ft)
		if err != nil {
			return nilValue, fmt.Errorf("%w while resolving argument %d for %s", err, i+1, funcInfo(fn))
		}
		in[i] = reflect.ValueOf(v)
	}

	out := fn.Call(in)

	if len(out) == 0 {
		return nilValue, nil
	}

	last := out[len(out)-1]
	if isError(last) {
		if err, ok := last.Interface().(error); ok {
			return nilValue, err
		}
	}
	return out[0], nil
}

func funcInfo(fn reflect.Value) string {
	buf := &bytes.Buffer{}
	t := fn.Type()
	fp := runtime.FuncForPC(fn.Pointer())
	if fp != nil {
		name := fp.Name()
		buf.WriteString(name)
		buf.WriteString("(")
		for i := 0; i < t.NumIn(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(typeName(t.In(i)))
		}
		buf.WriteString(")")

		if t.NumOut() > 0 {
			buf.WriteString(" ")
		}

		if t.NumOut() > 1 {
			buf.WriteString("(")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(typeName(t.Out(i)))
		}
		if t.NumOut() > 1 {
			buf.WriteString(")")
		}

		file, line := fp.FileLine(fn.Pointer())
		if file != "" {
			buf.WriteString(" @ (")
			buf.WriteString(file)
			buf.WriteString(":")
			buf.WriteString(fmt.Sprintf("%d", line))
			buf.WriteString(")")
		}
	}

	return buf.String()
}

func typeName(t reflect.Type) string {
	var name string
	if t.Kind() == reflect.Ptr {
		name = "*"
		t = t.Elem()
	}
	return name + t.Name()
}

func baseType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func isPtr(typ reflect.Type) bool {
	return typ.Kind() == reflect.Ptr
}

func isError(last reflect.Value) bool {
	return last.Type().Implements(errorInterface)
}

func validateProvider(t reflect.Type, provider any) {
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("provider must be a function: %+v", provider))
	}
	if t.NumOut() == 0 {
		panic(fmt.Sprintf("provider must return 1 value, or 1 value and an error but has no return types: %+v", provider))
	}
}
