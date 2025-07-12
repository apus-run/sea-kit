package slogx

import (
	"fmt"
)

// sprint is a function that takes a variadic parameter of type interface{} and returns a string.
// The function works as follows:
// - If no arguments are provided, it returns an empty string.
// - If a single argument is provided:
//   - If the argument is of type string, it returns the string as is.
//   - If the argument is not of type string but implements the fmt.Stringer interface, it returns the string representation of the argument.
//   - If the argument is not of type string and does not implement the fmt.Stringer interface, it converts the argument to a string using fmt.Sprint and returns the result.
//
// - If more than one argument is provided, it converts all arguments to a string using fmt.Sprint and returns the result.
func sprint(a ...any) string {
	if len(a) == 0 {
		return ""
	} else if len(a) == 1 {
		if s, ok := a[0].(string); ok {
			return s
		} else if v, ok := a[0].(fmt.Stringer); ok {
			return v.String()
		} else {
			return fmt.Sprint(a...)
		}
	} else {
		return fmt.Sprint(a...)
	}
}

// sprintf is a function that takes a string template and a variadic parameter of type interface{} and returns a string.
// The function works as follows:
// - If no arguments are provided, it returns the template string as is.
// - If the template string is not empty, it formats the string using fmt.Sprintf with the provided arguments and returns the result.
// - If only one argument is provided and it is of type string, it returns the string as is.
// - Otherwise, it converts the arguments to a string using the sprint function and returns the result.
func sprintf(template string, args ...any) string {
	if len(args) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, args...)
	}

	if len(args) == 1 {
		if str, ok := args[0].(string); ok {
			return str
		}
	}
	return sprint(args...)
}
