package inject

import (
	"fmt"
	"reflect"
)

func Invoke[T any](c Container, fn any, args ...any) (T, error) {
	var t T
	v, err := c.Invoke(fn, args...)
	if err != nil {
		return t, err
	}
	if tv, ok := v.(T); ok {
		t = tv
	} else {
		typ := reflect.TypeOf(t)
		return t, fmt.Errorf("unable convert return value to expected type from: %s got: %s %+v", funcInfo(reflect.ValueOf(fn)), typeName(typ), v)
	}
	return t, err
}

func MustInvoke[T any](c Container, fn any, args ...any) T {
	out, err := Invoke[T](c, fn, args...)
	if err != nil {
		panic(err)
	}
	return out
}

func Resolve[T any](c Container, typ T) (T, error) {
	c2, ok := c.(*container)
	var t T
	if !ok {
		return t, fmt.Errorf("argument is not *container: %+v", c)
	}
	out, err := c2.Resolve(typ)
	if tv, ok := out.(T); ok {
		t = tv
	} else {
		typ := reflect.TypeOf(t)
		return t, fmt.Errorf("unable convert value to expected type: %s got: %+v", typeName(typ), out)
	}
	return t, err
}

func MustResolve[T any](c Container, typ T) T {
	out, err := Resolve(c, typ)
	if err != nil {
		panic(err)
	}
	return out
}

func Decorate[T any](provider any, decorator func(T) T) any {
	baseFunc := reflect.ValueOf(provider)
	t := baseFunc.Type()
	validateProvider(t, provider)
	v := reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
		results = baseFunc.Call(args)
		// TODO check for error returns
		out, ok := results[0].Interface().(T)
		if ok {
			out = decorator(out)
			results[0] = reflect.ValueOf(out)
		}
		return
	})
	return v.Interface()
}

func Singleton(provider any) any {
	value := nilValue
	baseFunc := reflect.ValueOf(provider)
	t := baseFunc.Type()
	validateProvider(t, provider)
	v := reflect.MakeFunc(baseFunc.Type(), func(args []reflect.Value) (results []reflect.Value) {
		if value != nilValue {
			results = make([]reflect.Value, t.NumOut())
			results[0] = value
			return
		}
		results = baseFunc.Call(args)
		value = results[0]
		return
	})
	return v.Interface()
}
