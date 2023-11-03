package inject

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Container(t *testing.T) {
	type data struct {
		name string
	}

	type data1 struct {
		name string
	}

	type data2 struct {
		name string
	}

	d := &data{
		name: "1 name",
	}

	d1 := data1{
		name: "d1",
	}

	c := NewContainer()

	c.Bind(d, d1)

	c.Register(func(d data, d2 *data1) data2 {
		return data2{
			name: strings.ReplaceAll(d.name, "1", "2"),
		}
	})

	_, err := c.Invoke(func(d data2) {
		require.NotNil(t, d)
		require.Equal(t, "2 name", d.name)
	})

	require.NoError(t, err)

	type data3 struct {
		name string
	}

	type data4 struct {
		name string
	}

	c2 := NewContainer(c)

	c2.Register(func(d *data, d2 data2) *data3 {
		return &data3{
			name: d.name + ", " + d2.name,
		}
	})

	var d4 data4
	_, err = c2.Invoke(func(d data3) {
		d4 = data4{
			name: d.name,
		}
	})
	require.NoError(t, err)

	require.Equal(t, "1 name, 2 name", d4.name)

	c3 := NewContainer(c2)

	type err1 struct{}

	c3.Register(func() (err1, error) {
		return err1{}, fmt.Errorf("an error")
	})

	_, err = c3.Invoke(func(err1 err1) {})
	require.Error(t, err)

	var d1v2 data1
	_, err = c3.Invoke(func(d1 *data1) {
		d1v2 = *d1
	})
	require.NoError(t, err)
	require.Equal(t, "d1", d1v2.name)

	_, err = c3.Invoke(func(d1 *data1) error {
		return fmt.Errorf("direct error")
	})
	require.Error(t, err)

	c4 := NewContainer(c3)
	_, err = c4.Invoke(func(err1 err1) {})
	require.Error(t, err)

	dv3, err := c4.Resolve(data{})
	require.NoError(t, err)
	require.Equal(t, *d, dv3)

	dv4, err := c4.Resolve(&data{})
	require.NoError(t, err)
	require.Equal(t, d, dv4)
}

func Test_InjectingContainer(t *testing.T) {
	c := NewContainer()

	_, err := c.Invoke(func(c Container) {
		require.NotNil(t, c)
	})
	require.NoError(t, err)
}

func Test_UnresolvedDependency(t *testing.T) {
	c := NewContainer()

	type unbound struct{}

	_, err := c.Invoke(func(u unbound) {})

	require.Error(t, err)
	require.Contains(t, err.Error(), "unbound")
	require.Contains(t, err.Error(), "Test_UnresolvedDependency")
	require.Contains(t, err.Error(), "container_test.go")

	c.Bind(unbound{})
	_, err = c.Invoke(func(u unbound, u2 unbound, s string) {})

	require.Error(t, err)
	require.Contains(t, err.Error(), "string")
	require.Contains(t, err.Error(), "argument 3")
	require.Contains(t, err.Error(), "Test_UnresolvedDependency")
	require.Contains(t, err.Error(), "container_test.go")
}

func Test_ExtraArgs(t *testing.T) {
	c := NewContainer()

	type bound1 struct{}
	type bound2 struct{}

	c.Bind(bound1{}, bound2{})

	var got []any
	f := func(b1 bound1, b2 bound2, s1 string, s2 string, i1 int) {
		got = append(got, s1, s2, i1)
	}

	_, err := c.Invoke(f, "the s1", "the s2", 95)

	require.NoError(t, err)
	require.Equal(t, []any{"the s1", "the s2", 95}, got)
}
