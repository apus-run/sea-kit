package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	sec := time.Now().Unix()
	tz := "Asia/Shanghai"
	actual := Format(sec, "YYYY-MM-DD HH:mm:ss", tz)
	t.Logf("actual: %s", actual)
	_, _ = time.LoadLocation(tz)
	expected := time.Unix(sec, 0).Format("2006-01-02 15:04:05")
	assert.Equal(t, expected, actual)
}
