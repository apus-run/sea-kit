package stringx

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringRing(t *testing.T) {
	r := NewStringRing(3)
	t.Logf("[]string{} =  %v", r.GetAll())
	r.Put("a")
	t.Logf(`[]string{"a"} =  %v`, r.GetAll())
	r.Put("b")
	t.Logf(`[]string{"a", "b"} = %v`, r.GetAll())
	r.Put("c")
	t.Logf(`[]string{"a", "b", "c"} = %v`, r.GetAll())
	r.Put("d")
	t.Logf(`[]string{"b", "c", "d"} = %v`, r.GetAll())
}

func TestStringRingConcurrent(t *testing.T) {

	// Check for racy behavior by spinning up a bunch of goroutines to write
	// to the ring and verifying that it ends up with the correct set of
	// entries.
	r := NewStringRing(1000)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				r.Put(fmt.Sprintf("%d:%d", i, j))
			}
		}(i)
	}
	wg.Wait()
	require.Equal(t, 1000, len(r.GetAll()))
}
