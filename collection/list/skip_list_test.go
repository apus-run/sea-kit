package list

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apus-run/sea-kit/collection/list/internal"
)

func TestNewSkipList(t *testing.T) {
	testCases := []struct {
		name      string
		compare   internal.Comparator[int]
		wantSlice []int
	}{
		{
			name:      "new skip list",
			compare:   internal.ComparatorRealNumber[int],
			wantSlice: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList(tc.compare)
			assert.Equal(t, tc.wantSlice, sl.AsSlice())
		})
	}
}

func TestSkipList_AsSlice(t *testing.T) {
	testCases := []struct {
		name      string
		compare   internal.Comparator[int]
		wantSlice []int
	}{
		{
			name:      "no err is ok",
			compare:   internal.ComparatorRealNumber[int],
			wantSlice: []int{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			assert.Equal(t, tc.wantSlice, sl.AsSlice())
		})
	}
}

func TestSkipList_Cap(t *testing.T) {
	testCases := []struct {
		name     string
		compare  internal.Comparator[int]
		wantSize int
	}{
		{
			name:     "no err is ok",
			compare:  internal.ComparatorRealNumber[int],
			wantSize: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			assert.Equal(t, tc.wantSize, sl.Cap())
		})
	}
}

func TestSkipList_DeleteElement(t *testing.T) {
	testCases := []struct {
		name     string
		compare  internal.Comparator[int]
		value    int
		wantBool bool
	}{
		{
			name:     "no err is ok",
			compare:  internal.ComparatorRealNumber[int],
			value:    1,
			wantBool: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			ok := sl.DeleteElement(tc.value)
			assert.Equal(t, tc.wantBool, ok)
		})
	}
}

func TestSkipList_Insert(t *testing.T) {
	testCases := []struct {
		name      string
		compare   internal.Comparator[int]
		key       int
		wantSlice []int
	}{
		{
			name:      "no err is ok",
			compare:   internal.ComparatorRealNumber[int],
			key:       1,
			wantSlice: []int{1},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			sl.Insert(tc.key)
			assert.Equal(t, tc.wantSlice, sl.AsSlice())
		})
	}
}

func TestSkipList_Len(t *testing.T) {
	testCases := []struct {
		name     string
		compare  internal.Comparator[int]
		wantSize int
	}{
		{
			name:     "no err is ok",
			compare:  internal.ComparatorRealNumber[int],
			wantSize: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			assert.Equal(t, tc.wantSize, sl.Len())
		})
	}
}

func TestSkipList_Search(t *testing.T) {
	testCases := []struct {
		name     string
		compare  internal.Comparator[int]
		value    int
		wantBool bool
	}{
		{
			name:     "no err is ok",
			compare:  internal.ComparatorRealNumber[int],
			value:    1,
			wantBool: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			ok := sl.Search(tc.value)
			assert.Equal(t, tc.wantBool, ok)
		})
	}
}
