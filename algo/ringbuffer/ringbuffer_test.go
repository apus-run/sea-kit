package ringbuffer

import (
	"reflect"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	const numItems = 10
	rb := New[int](numItems)

	for i := 0; i < numItems-1; i++ {
		rb.Add(i)
	}

	t.Run("NotFull", func(t *testing.T) {
		if ll := rb.Len(); ll != numItems-1 {
			t.Fatalf("got len %d; want %d", ll, numItems-1)
		}
		all := rb.Take()
		want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
		if !reflect.DeepEqual(all, want) {
			t.Fatalf("items mismatch\ngot: %v\nwant %v", all, want)
		}
	})

	t.Run("Full", func(t *testing.T) {
		// Append items to evict something
		rb.Add(98)
		rb.Add(99)

		if ll := rb.Len(); ll != numItems {
			t.Fatalf("got len %d; want %d", ll, numItems)
		}
		all := rb.Take()
		want := []int{1, 2, 3, 4, 5, 6, 7, 8, 98, 99}
		if !reflect.DeepEqual(all, want) {
			t.Fatalf("items mismatch\ngot: %v\nwant %v", all, want)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		rb.Clear()
		if ll := rb.Len(); ll != 0 {
			t.Fatalf("got len %d; want 0", ll)
		}
		all := rb.Take()
		if len(all) != 0 {
			t.Fatalf("got non-empty list; want empty")
		}
	})
}
