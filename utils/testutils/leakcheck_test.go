package testutils

import (
	"testing"
)

func TestVerifyGoLeaksOnce(t *testing.T) {
	defer VerifyGoLeaksOnce(t)
}

func TestMain(m *testing.M) {
	VerifyGoLeaks(m)
}
