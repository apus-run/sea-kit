package do

import (
	"errors"
	"testing"
)

func TestDo_err(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")
	errs := Do(
		func() error { return a },
		func() error { return b },
	)
	if !errors.Is(errs, a) {
		t.Fatal(errs)
	}
	if !errors.Is(errs, b) {
		t.Fatal(errs)
	}
}
