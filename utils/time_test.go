package utils

import (
	"context"
	"testing"
	"time"
)

func TestUnixMill(t *testing.T) {
	tt := time.UnixMilli(1655038017807)
	t.Log(tt)
}

func TestFormatISO(t *testing.T) {
	t.Log(time.Now().Format(time.RFC3339))
}

func TestShrink(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c := context.Background()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) != time.Second {
		t.Fatalf("new timeout must be equal 1 second")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Second || time.Until(deadline) < time.Millisecond*500 {
		t.Fatalf("ctx deadline must be less than 1s and greater than 500ms")
	}
}

func TestShrinkWithTimeout(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) != time.Second {
		t.Fatalf("new timeout must be equal 1 second")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Second || time.Until(deadline) < time.Millisecond*500 {
		t.Fatalf("ctx deadline must be less than 1s and greater than 500ms")
	}
}

func TestShrinkWithDeadline(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) >= time.Millisecond*500 {
		t.Fatalf("new timeout must be less than 500 ms")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Millisecond*500 || time.Until(deadline) < time.Millisecond*200 {
		t.Fatalf("ctx deadline must be less than 500ms and greater than 200ms")
	}
}

func TestCurrTimeGet(t *testing.T) {
	curr := CurrTimeGet()
	println(curr)
}

func TestUnixTimeToDateTime(t *testing.T) {
	dateTime := UnixTimeToDateTime(1419238997)
	if dateTime != 20141222170317 {
		t.Errorf("err in UnixTimeToDateTime(), ok:20141222170317, now:%d", dateTime)
	}
}

func TestTimestampSplit(t *testing.T) {
	// good case
	timestr := "20160808144120"
	date, time, err := TimestampSplit(timestr)
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}
	if date != "20160808" {
		t.Errorf("unexpected split result(date): %s", date)
	}
	if time != "144120" {
		t.Errorf("unexpected split result(time): %s", time)
	}

	// bad case 1:
	timestr = "2016080814412011111"
	_, _, err = TimestampSplit(timestr)
	if err == nil {
		t.Errorf("err should happen for : %s", timestr)
	}

	// bad case 2:
	timestr = "20160808144"
	_, _, err = TimestampSplit(timestr)
	if err == nil {
		t.Errorf("err should happen for : %s", timestr)
	}
}

func TestStrToUnix(t *testing.T) {
	ts, err := StrToUnix("2017-05-10 19:35:41")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
	if ts != 1494416141 {
		t.Errorf("wrong timestamp")
		return
	}
}

func TestUnixToStr(t *testing.T) {
	if UnixToStr(1494416141) != "2017-05-10 19:35:41" {
		t.Errorf("wrong time string")
		return
	}
}

func TestWaitTill_case1(t *testing.T) {
	waitSecs := int64(2)
	start := time.Now().Unix()

	toTime := start + waitSecs

	WaitTill(toTime)

	passSecs := time.Now().Unix() - start

	if passSecs != waitSecs {
		t.Errorf("err in WaitTill(): wait=%d, pass=%d", waitSecs, passSecs)
	}
}

func TestWaitTill_case2(t *testing.T) {
	start := time.Now().Unix()

	toTime := start - 2

	WaitTill(toTime)

	passSecs := time.Now().Unix() - start

	if passSecs != 0 {
		t.Errorf("err in WaitTill(): wait=0, pass=%d", passSecs)
	}
}
