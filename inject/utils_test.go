package inject

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Decorate(t *testing.T) {
	type t1 struct {
		count int
	}
	c := NewContainer()

	inst := t1{
		count: 0,
	}
	c.Register(Decorate(func() t1 {
		inst.count = 1
		return inst
	}, func(a t1) t1 {
		a.count = 100
		return a
	}))

	_, err := c.Invoke(func(t1 t1) {
		require.Equal(t, 100, t1.count)
	})
	require.NoError(t, err)
}

func Test_NonSingleton(t *testing.T) {
	type t1 struct {
		count int
	}
	c := NewContainer()

	invocationCount := 0
	inst := t1{
		count: 0,
	}
	c.Register(func() t1 {
		invocationCount++
		inst.count++
		return inst
	})

	_, err := c.Invoke(func(t1 t1) {
		require.Equal(t, 1, t1.count)
	})
	require.NoError(t, err)
	require.Equal(t, 1, invocationCount)

	_, err = c.Invoke(func(t1 t1) {
		require.Equal(t, 2, t1.count)
	})
	require.NoError(t, err)
	require.Equal(t, 2, invocationCount)
}

func Test_Singleton(t *testing.T) {
	type t1 struct {
		count int
	}
	c := NewContainer()

	invocationCount := 0
	inst := t1{
		count: 0,
	}
	c.Register(Singleton(func() t1 {
		invocationCount++
		inst.count++
		return inst
	}))

	_, err := c.Invoke(func(t1 t1) {
		require.Equal(t, 1, t1.count)
	})
	require.NoError(t, err)
	require.Equal(t, 1, invocationCount)

	_, err = c.Invoke(func(t1 t1) {
		require.Equal(t, 1, t1.count)
	})
	require.NoError(t, err)
	require.Equal(t, 1, invocationCount)
}

func Test_Resolve(t *testing.T) {
	c := NewContainer()

	type t1 struct {
		v1 int
	}

	type t2 struct {
		v2 int
	}

	c.Bind(t1{v1: 99})

	v, err := Resolve(c, t1{})
	require.NoError(t, err)
	require.Equal(t, 99, v.v1)

	_, err = Resolve(c, t2{})
	require.Error(t, err)
}

func Test_MustResolve(t *testing.T) {
	c := NewContainer()

	type t1 struct {
		v1 int
	}

	type t2 struct {
		v2 int
	}

	c.Bind(t1{v1: 99})

	v := MustResolve(c, t1{})
	require.Equal(t, 99, v.v1)

	require.Panics(t, func() {
		_ = MustResolve(c, t2{})
	})
}
