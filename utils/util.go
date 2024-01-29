package utils

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	timeUnixZero = time.Unix(0, 0).UTC()
)

// In returns true if |s| is *in* |a| slice.
func In(s string, a []string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

// ContainsAny returns true if |s| contains any element of |a|.
func ContainsAny(s string, a []string) bool {
	for _, x := range a {
		if strings.Contains(s, x) {
			return true
		}
	}
	return false
}

// Index returns the index of |s| *in* |a| slice, and -1 if not found.
func Index(s string, a []string) int {
	for i, x := range a {
		if x == s {
			return i
		}
	}
	return -1
}

// SSliceEqual returns true if the given string slices are equal
func SSliceEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i, aa := range a {
		if aa != b[i] {
			return false
		}
	}
	return true
}

// Reverse returns the given slice of strings in reverse order.
func Reverse(s []string) []string {
	r := make([]string, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		r = append(r, s[i])
	}
	return r
}

// insertString inserts the given string into the slice at the given index.
func insertString(strs []string, idx int, s string) []string {
	oldLen := len(strs)
	strs = append(strs, "")
	copy(strs[idx+1:], strs[idx:oldLen])
	strs[idx] = s
	return strs
}

// InsertStringSorted inserts the given string into the sorted slice of strings
// if it does not already exist. Maintains sorted order.
func InsertStringSorted(strs []string, s string) []string {
	idx := sort.SearchStrings(strs, s)
	if idx == len(strs) || strs[idx] != s {
		return insertString(strs, idx, s)
	}
	return strs
}

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// MaxInt returns the largest integer of the arguments provided.
func MaxInt(intList ...int) int {
	ret := intList[0]
	for _, i := range intList[1:] {
		if i > ret {
			ret = i
		}
	}
	return ret
}

// MaxInt64 returns largest integer of a and b.
func MaxInt64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

// MaxInt32 returns largest integer of a and b.
func MaxInt32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

// MinInt returns the smaller integer of a and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MinInt64 returns the smaller integer of a and b.
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// MinInt32 returns the smaller integer of a and b.
func MinInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// AbsInt returns the absolute value of v.
func AbsInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// TimeStampMs returns the current time in milliseconds since the epoch.
func TimeStampMs() int64 {
	return TimeStamp(time.Millisecond)
}

// TimeStamp returns the current time in the units defined by the given target unit.
// e.g. TimeStamp(time.Millisecond) will return the time in Milliseconds.
// The result is always rounded down to the lowest integer from the
// representation in nano seconds.
func TimeStamp(targetUnit time.Duration) int64 {
	return time.Now().UnixNano() / int64(targetUnit)
}

// RepeatJoin repeats a given string N times with the given separator between
// each instance.
func RepeatJoin(str, sep string, n int) string {
	if n <= 0 {
		return ""
	}
	return str + strings.Repeat(sep+str, n-1)
}

// AddParams adds the second instance of safemap[string]string to the first and
// returns the first safemap.
func AddParams(a map[string]string, b ...map[string]string) map[string]string {
	if a == nil {
		a = make(map[string]string, len(b))
	}
	for _, oneMap := range b {
		for k, v := range oneMap {
			a[k] = v
		}
	}
	return a
}

// CopyStringMap returns a copy of the provided safemap[string]string such that
// reflect.DeepEqual returns true for the given safemap and the returned safemap. In
// particular, preserves nil input.
func CopyStringMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	ret := make(map[string]string, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

// CopyStringSlice copies the given []string such that reflect.DeepEqual returns
// true for the given slice and the returned slice. In particular, preserves
// nil slice input.
func CopyStringSlice(s []string) []string {
	if s == nil {
		return nil
	}
	rv := make([]string, len(s))
	copy(rv, s)
	return rv
}

// CopyString returns a copy of the given string. This may seem unnecessary, but
// is very important at preventing leaks of strings. For example, subslicing
// a string can prevent the larger string from being cleaned up.
func CopyString(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := &strings.Builder{}
	b.WriteString(s)
	return b.String()
}

// IsNil returns true if i is nil or is an interface containing a nil or invalid value.
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		inner := v.Elem()
		if !inner.IsValid() {
			return true
		}
		if inner.CanInterface() {
			return IsNil(inner.Interface())
		}
		return false
	default:
		return false
	}
}

// Repeat calls the provided function 'fn' immediately and then in intervals
// defined by 'interval'. If anything is sent on the provided stop channel,
// the iteration stops.
func Repeat(interval time.Duration, stopCh <-chan bool, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	fn()
MainLoop:
	for {
		select {
		case <-stopCh:
			break MainLoop
		case <-ticker.C:
			fn()
		}
	}
}

// RepeatCtx calls the provided function 'fn' immediately and then in intervals
// defined by 'interval'. If the given context is canceled, the iteration stops.
func RepeatCtx(ctx context.Context, interval time.Duration, fn func(ctx context.Context)) {
	if interval <= 0 {
		return
	}
	ticker := time.NewTicker(interval)
	done := ctx.Done()
	defer ticker.Stop()
	fn(ctx)
MainLoop:
	for {
		select {
		case <-done:
			break MainLoop
		case <-ticker.C:
			fn(ctx)
		}
	}
}

// Truncate the given string to the given length. If the string was shortened,
// change the last three characters to ellipses, unless the specified length is
// 3 or less.
func Truncate(s string, length int) string {
	if len(s) > length {
		if length <= 3 {
			return s[:length]
		}
		ellipses := "..."
		return s[:length-len(ellipses)] + ellipses
	}
	return s
}

// TruncateNoEllipses truncates the given string to the given length, without
// the use of ellipses.
func TruncateNoEllipses(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	return s
}

// FirstNonEmpty returns the first of its args that is not "". It is useful when a certain value
// would be preferred if present but others are available as fallbacks. If all of its args are "",
// returns "".
func FirstNonEmpty(args ...string) string {
	a := ""
	for _, a = range args {
		if a != "" {
			return a
		}
	}
	return a
}

// SplitLines returns a slice of the lines of s, split on newline characters. If the input string
// ends in a single newline, we strip it rather than returning a blank extra line.
//
// Note that this currently works only with UNIX-style line breaks, not DOS \r\n ones, though
// support for those may be added later.
func SplitLines(s string) []string {
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n")
}

// ParseIntSet parses a string expression like "5", "3-8", or "3,4,9" into a
// slice of integers: [5], [3, 4, 5, 6, 7, 8], [3, 4, 9].
func ParseIntSet(expr string) ([]int, error) {
	rv := []int{}
	if expr == "" {
		return rv, nil
	}
	ranges := strings.Split(expr, ",")
	for _, r := range ranges {
		endpoints := strings.Split(r, "-")
		if len(endpoints) == 1 {
			v, err := strconv.Atoi(endpoints[0])
			if err != nil {
				return nil, err
			}
			rv = append(rv, v)
		} else if len(endpoints) == 2 {
			if endpoints[0] == "" {
				return nil, fmt.Errorf("Invalid expression %q", r)
			}
			start, err := strconv.Atoi(endpoints[0])
			if err != nil {
				return nil, err
			}
			if endpoints[1] == "" {
				return nil, fmt.Errorf("Invalid expression %q", r)
			}
			end, err := strconv.Atoi(endpoints[1])
			if err != nil {
				return nil, err
			}
			if start > end {
				return nil, fmt.Errorf("Cannot have a range whose beginning is greater than its end (%d vs %d)", start, end)
			}
			for i := start; i <= end; i++ {
				rv = append(rv, i)
			}
		} else {
			return nil, fmt.Errorf("Invalid expression %q", r)
		}
	}
	return rv, nil
}
