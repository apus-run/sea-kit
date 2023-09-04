package utils

import (
	"reflect"
	"testing"
)

func TestChunkSlice(t *testing.T) {
	type args struct {
		slice     []string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "slice with elements",
			args: args{
				slice:     []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
				chunkSize: 3,
			},
			want: [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
				{"0"},
			},
		},
		{
			name: "slice with only one element",
			args: args{
				slice:     []string{"1"},
				chunkSize: 3,
			},
			want: [][]string{
				{"1"},
			},
		},
		{
			name: "empty slice",
			args: args{
				slice:     []string{},
				chunkSize: 3,
			},
			want: make([][]string, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChunkSlice(tt.args.slice, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChunkSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
