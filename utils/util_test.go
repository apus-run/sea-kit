package utils

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestFirstNonEmpty(t *testing.T) {
	assert.Equal(t, "", FirstNonEmpty())
	assert.Equal(t, "", FirstNonEmpty(""))
	assert.Equal(t, "a", FirstNonEmpty("a", "b"))
	assert.Equal(t, "c", FirstNonEmpty("", "", "c"))
}

func TestSplitLines_StripsTrailingNewline(t *testing.T) {
	assert.Equal(t, []string{"this", "that"}, SplitLines("this\nthat\n"))
}

func TestParseIntSet(t *testing.T) {

	test := func(input string, expect []int, expectErr string) {
		res, err := ParseIntSet(input)
		if expectErr != "" {
			require.Contains(t, err.Error(), expectErr)
		} else {
			require.NoError(t, err)
			require.Equal(t, expect, res)
		}
	}
	test("", []int{}, "")
	test("19", []int{19}, "")
	test("1,2,3", []int{1, 2, 3}, "")
	test("1-3", []int{1, 2, 3}, "")
	test("1,2,4-6", []int{1, 2, 4, 5, 6}, "")
	test("a", nil, "parsing \"a\": invalid syntax")
	test(" 4, 6, 9 - 11", nil, "parsing \" 4\": invalid syntax")
	test("4-9-10", nil, "Invalid expression \"4-9-10\"")
	test("9-3", nil, "Cannot have a range whose beginning is greater than its end (9 vs 3)")
	test("1-3,11-13,21-23", []int{1, 2, 3, 11, 12, 13, 21, 22, 23}, "")
	test("-2", nil, "Invalid expression \"-2\"")
	test("2-", nil, "Invalid expression \"2-\"")
}

func TestTruncate(t *testing.T) {
	s := "abcdefghijkl"
	require.Equal(t, "", Truncate(s, 0))
	require.Equal(t, "a", Truncate(s, 1))
	require.Equal(t, "ab", Truncate(s, 2))
	require.Equal(t, "abc", Truncate(s, 3))
	require.Equal(t, "a...", Truncate(s, 4))
	require.Equal(t, "ab...", Truncate(s, 5))
	require.Equal(t, s, Truncate(s, len(s)))
	require.Equal(t, s, Truncate(s, len(s)+1))
}

func TestInsertStringSorted(t *testing.T) {
	require.Equal(t, []string{"a"}, InsertStringSorted([]string{}, "a"))
	require.Equal(t, []string{"a"}, InsertStringSorted([]string{"a"}, "a"))
	require.Equal(t, []string{"a", "b"}, InsertStringSorted([]string{"a"}, "b"))
	require.Equal(t, []string{"0", "a", "b"}, InsertStringSorted([]string{"a", "b"}, "0"))
	require.Equal(t, []string{"0", "a", "b"}, InsertStringSorted([]string{"0", "a", "b"}, "b"))
}

func TestIsNil(t *testing.T) {
	require.True(t, IsNil(nil))
	require.False(t, IsNil(false))
	require.False(t, IsNil(0))
	require.False(t, IsNil(""))
	require.False(t, IsNil([0]int{}))
	type Empty struct{}
	require.False(t, IsNil(Empty{}))
	require.True(t, IsNil(chan interface{}(nil)))
	require.False(t, IsNil(make(chan interface{})))
	var f func()
	require.True(t, IsNil(f))
	require.False(t, IsNil(func() {}))
	require.True(t, IsNil(map[bool]bool(nil)))
	require.False(t, IsNil(make(map[bool]bool)))
	require.True(t, IsNil([]int(nil)))
	require.False(t, IsNil([][]int{nil}))
	require.True(t, IsNil((*int)(nil)))
	var i int
	require.False(t, IsNil(&i))
	var pi *int
	require.True(t, IsNil(pi))
	require.True(t, IsNil(&pi))
	var ppi **int
	require.True(t, IsNil(&ppi))
	var c chan interface{}
	require.True(t, IsNil(&c))
	var w io.Writer
	require.True(t, IsNil(w))
	w = (*bytes.Buffer)(nil)
	require.True(t, IsNil(w))
	w = &bytes.Buffer{}
	require.False(t, IsNil(w))
	require.False(t, IsNil(&w))
	var ii interface{}
	ii = &pi
	require.True(t, IsNil(ii))
}

func TestInsertString(t *testing.T) {
	require.Equal(t, []string{"a"}, insertString([]string{}, 0, "a"))
	require.Equal(t, []string{"b", "a"}, insertString([]string{"a"}, 0, "b"))
	require.Equal(t, []string{"b", "c", "a"}, insertString([]string{"b", "a"}, 1, "c"))
	require.Equal(t, []string{"b", "c", "a", "d"}, insertString([]string{"b", "c", "a"}, 3, "d"))
}
