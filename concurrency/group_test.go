package concurrency

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupGet(t *testing.T) {
	count := 0
	g := NewGroup(func() interface{} {
		count++
		return count
	})
	v := g.Get("/x/internal/dummy/user")
	assert.Equal(t, 1, v.(int))

	v = g.Get("/x/internal/dummy/avatar")
	assert.Equal(t, 2, v.(int))

	v = g.Get("/x/internal/dummy/user")
	assert.Equal(t, 1, v.(int))
	assert.Equal(t, 2, count)
}

func TestGroupReset(t *testing.T) {
	g := NewGroup(func() interface{} {
		return 1
	})
	g.Get("/x/internal/dummy/user")
	call := false
	g.Reset(func() interface{} {
		call = true
		return 1
	})

	length := 0
	for range g.objs {
		length++
	}

	assert.Equal(t, 0, length)

	g.Get("/x/internal/dummy/user")
	assert.Equal(t, true, call)
}

func TestGroupClear(t *testing.T) {
	g := NewGroup(func() interface{} {
		return 1
	})
	g.Get("/x/internal/dummy/user")
	length := 0
	for range g.objs {
		length++
	}
	assert.Equal(t, 1, length)

	g.Clear()
	length = 0
	for range g.objs {
		length++
	}
	assert.Equal(t, 0, length)
}

type Counter struct {
	Value int
}

func (c *Counter) Incr() {
	c.Value++
}

func ExampleGroup_Get() {
	new := func() interface{} {
		fmt.Println("Only Once")
		return &Counter{}
	}
	group := NewGroup(new)

	// Create a new Counter
	group.Get("pass").(*Counter).Incr()

	// Get the created Counter again.
	group.Get("pass").(*Counter).Incr()
	// Output:
	// Only Once
}

func ExampleGroup_Reset() {
	new := func() interface{} {
		return &Counter{}
	}
	group := NewGroup(new)

	newV2 := func() interface{} {
		fmt.Println("New V2")
		return &Counter{}
	}
	// Reset the new function and clear all created objects.
	group.Reset(newV2)

	// Create a new Counter
	group.Get("pass").(*Counter).Incr()
	// Output:
	// New V2
}
