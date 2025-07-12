package slogx

import (
	"fmt"
	"testing"
)

func Test_sprint(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{
			name: "NoArgs",
			args: []interface{}{},
			want: "",
		},
		{
			name: "WithOneArgString",
			args: []interface{}{"arg1"},
			want: "arg1",
		},
		{
			name: "WithOneArgNotString",
			args: []interface{}{123},
			want: "123",
		},
		{
			name: "WithMultipleArgsString",
			args: []interface{}{"arg1", "arg2"},
			want: "arg1arg2",
		},
		{
			name: "WithMultipleArgsNotString",
			args: []interface{}{123, 456},
			want: "123 456",
		},
		{
			name: "WithErrorArgs",
			args: []interface{}{fmt.Errorf("error message")},
			want: "error message",
		},
		{
			name: "WithStringerArgs",
			args: []interface{}{stringer{str: "stringer"}},
			want: "stringer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sprint(tt.args...); got != tt.want {
				t.Errorf("sprint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sprintf(t *testing.T) {
	type args struct {
		template string
		args     []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "NoArgs",
			args: args{template: "template", args: []interface{}{}},
			want: "template",
		},
		{
			name: "WithTemplateAndOneArg",
			args: args{template: "template %s", args: []interface{}{"arg1"}},
			want: "template arg1",
		},
		{
			name: "WithTemplateAndMultipleArgs",
			args: args{template: "template %s %s", args: []interface{}{"arg1", "arg2"}},
			want: "template arg1 arg2",
		},
		{
			name: "WithOneArgNotString",
			args: args{template: "", args: []interface{}{123}},
			want: "123",
		},
		{
			name: "WithMultipleArgsNotString",
			args: args{template: "", args: []interface{}{123, 456}},
			want: "123 456",
		},
		{
			name: "WithErrorArgs",
			args: args{template: "", args: []interface{}{fmt.Errorf("error message")}},
			want: "error message",
		},
		{
			name: "WithStringerArgs",
			args: args{template: "", args: []interface{}{stringer{str: "stringer"}}},
			want: "stringer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sprintf(tt.args.template, tt.args.args...); got != tt.want {
				t.Errorf("sprintf() = %v, want %v", got, tt.want)
			}
		})
	}
}

type stringer struct {
	str string
}

func (s stringer) String() string {
	return s.str
}
