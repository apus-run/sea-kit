package utils

import (
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
