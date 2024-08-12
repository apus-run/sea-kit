package collection

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		i := 1
		p := ToPtr[int](i)
		assert.Equal(t, &i, p)
	})
	t.Run("bool", func(t *testing.T) {
		i := true
		p := ToPtr[bool](i)
		assert.Equal(t, &i, p)
	})

	t.Run("string", func(t *testing.T) {
		s := "hello"
		p := ToPtr[string](s)
		assert.Equal(t, &s, p)
	})
}

func TestAllPtrFieldsNil(t *testing.T) {
	testCases := []struct {
		obj      interface{}
		expected bool
	}{
		{struct{}{}, true},
		{struct{ Foo int }{12345}, true},
		{&struct{ Foo int }{12345}, true},
		{struct{ Foo *int }{nil}, true},
		{&struct{ Foo *int }{nil}, true},
		{struct {
			Foo int
			Bar *int
		}{12345, nil}, true},
		{&struct {
			Foo int
			Bar *int
		}{12345, nil}, true},
		{struct {
			Foo *int
			Bar *int
		}{nil, nil}, true},
		{&struct {
			Foo *int
			Bar *int
		}{nil, nil}, true},
		{struct{ Foo *int }{new(int)}, false},
		{&struct{ Foo *int }{new(int)}, false},
		{struct {
			Foo *int
			Bar *int
		}{nil, new(int)}, false},
		{&struct {
			Foo *int
			Bar *int
		}{nil, new(int)}, false},
		{(*struct{})(nil), true},
	}
	for i, tc := range testCases {
		name := fmt.Sprintf("case[%d]", i)
		t.Run(name, func(t *testing.T) {
			if actual := AllPtrFieldsNil(tc.obj); actual != tc.expected {
				t.Errorf("%s: expected %t, got %t", name, tc.expected, actual)
			}
		})
	}
}

func TestRef(t *testing.T) {
	type T int

	val := T(0)
	pointer := ToPtr(val)
	if *pointer != val {
		t.Errorf("expected %d, got %d", val, *pointer)
	}

	val = T(1)
	pointer = ToPtr(val)
	if *pointer != val {
		t.Errorf("expected %d, got %d", val, *pointer)
	}
}

func TestFrom(t *testing.T) {
	assert.Equal(t, 543, From(ToPtr(543)))
	assert.Equal(t, "Alice", From(ToPtr("Alice")))
	assert.Zero(t, From[int](nil))
	assert.Nil(t, From[interface{}](nil))
	assert.Nil(t, From(ToPtr[fmt.Stringer](nil)))
}

func TestFromOr(t *testing.T) {
	type T int

	var val, def T = 1, 0

	out := FromOr(&val, def)
	if out != val {
		t.Errorf("expected %d, got %d", val, out)
	}

	out = FromOr(nil, def)
	if out != def {
		t.Errorf("expected %d, got %d", def, out)
	}
}

func TestIsNil(t *testing.T) {
	assert.False(t, IsNil(ToPtr(1)))
	assert.True(t, IsNil[int](nil))
}

func TestClone(t *testing.T) {
	assert.True(t, IsNil(Clone(((*int)(nil)))))

	v := 1
	assert.True(t, Clone(&v) != &v)
	assert.True(t, Equal(Clone(&v), &v))

	src := ToPtr(1)
	dst := Clone(&src)
	assert.Equal(t, &src, dst)
	assert.True(t, src == *dst)
}

func TestCloneBy(t *testing.T) {
	assert.True(t, IsNil(CloneBy(((**int)(nil)), Clone[int])))

	src := ToPtr(1)
	dst := CloneBy(&src, Clone[int])
	assert.Equal(t, &src, dst)
	assert.False(t, src == *dst)
}

func TestEqual(t *testing.T) {
	type T int

	if !Equal[T](nil, nil) {
		t.Errorf("expected true (nil == nil)")
	}
	if !Equal(ToPtr(T(123)), ToPtr(T(123))) {
		t.Errorf("expected true (val == val)")
	}
	if Equal(nil, ToPtr(T(123))) {
		t.Errorf("expected false (nil != val)")
	}
	if Equal(ToPtr(T(123)), nil) {
		t.Errorf("expected false (val != nil)")
	}
	if Equal(ToPtr(T(123)), ToPtr(T(456))) {
		t.Errorf("expected false (val != val)")
	}
}

func TestEqualTo(t *testing.T) {
	assert.True(t, EqualTo(ToPtr(1), 1))
	assert.False(t, EqualTo(ToPtr(2), 1))
	assert.False(t, EqualTo(nil, 0))
}

func TestMap(t *testing.T) {
	i := 1
	assert.Equal(t, ToPtr("1"), Map(&i, strconv.Itoa))
	assert.True(t, Map(nil, strconv.Itoa) == nil)

	assert.NotPanics(t, func() {
		_ = Map(nil, func(int) string {
			panic("Q_Q")
		})
	})

	assert.Panics(t, func() {
		_ = Map(&i, func(int) string {
			panic("Q_Q")
		})
	})
}
