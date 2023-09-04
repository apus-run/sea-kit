package utils

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"
	xtime "time"
)

const (
	YYYYMMDDHHmmssNoSplit = "20060102150405"
	YYYYMMDD              = "2006-01-02"
	YYYYMMDDHHmmSS        = "2006-01-02 15:04:05"
	YMDHms                = "2006.01.02 15:04:05"

	YYYYMMDDHHmmSSZone    = "2006-01-02 15:04:04 -0700"
	YYYYMMDDHHmmSSISO8601 = "2006-01-02T15:04:05.000Z"
)

// Time be used to MySql timestamp converting.
type Time int64

// Scan scan time.
func (jt *Time) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case xtime.Time:
		*jt = Time(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = Time(i)
	}
	return
}

// Value get time value.
func (jt Time) Value() (driver.Value, error) {
	return xtime.Unix(int64(jt), 0), nil
}

// Time get time.
func (jt Time) Time() xtime.Time {
	return xtime.Unix(int64(jt), 0)
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration xtime.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := xtime.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := xtime.Until(deadline); ctimeout < xtime.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, xtime.Duration(d))
	return d, ctx, cancel
}

func StringToTime(myTime, format string) (*xtime.Time, error) {
	parse, err := xtime.ParseInLocation(myTime, format, xtime.Local)
	if err != nil {
		return nil, err
	}
	return &parse, nil
}

func TimeToString(myTime xtime.Time, format string) string {
	return myTime.Format(format)
}

func TimeParseInt64(tstamp int64) *xtime.Time {
	if tstamp == 0 {
		return nil
	}
	t := xtime.Unix(tstamp, 0)
	return &t
}
func TimeUnixShowLayoutString(t int64) string {
	return TimeParseInt64(t).Format(YYYYMMDDHHmmSS)
}
func TimeShowLayoutString(t xtime.Time) string {
	return t.Format(YYYYMMDDHHmmSS)
}

func TimePushDayToString(t xtime.Time, day int) string {
	return t.AddDate(0, 0, day).Format(YYYYMMDDHHmmSS)
}

func DateTimeToDate(t xtime.Time) xtime.Time {
	return xtime.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, xtime.Local)
}

// CurrTimeGet gets current time in format like 2006-01-02 15:04:05
func CurrTimeGet() string {
	return xtime.Now().Format("2006-01-02 15:04:05")
}

// UnixTimeToDateTime converts from unix time to date-and-time
// e.g., from 13923223442(int64) to 20140611120132(int64)
//
// Params:
//   - unixTime: the number of seconds elapsed since January 1, 1970 UTC.
//
// Returns:
//
//	(datetime, error)
//	datetime - yyyymmddhhmmss, e.g., 20140611120132 (int64)
func UnixTimeToDateTime(unixTime int64) int64 {
	// convert from unix time to string of date-and-time
	timeStr := xtime.Unix(unixTime, 0).Format("20060102150405")

	// convert date-and-time from string to int
	timeInt, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		// this should not happen
		return 0
	}

	return timeInt
}

// TimestampSplit splits a full time string in format "yyyyMMddHHmmss"(e.g., "20130411143020"), to
// "yyyyMMdd" and "HHmmss"
//
// Params:
//   - timestamp: time string in format "yyyyMMddHHmmss"
//
// Returns:
//
//	("yyyyMMdd", "HHmmss", error)
func TimestampSplit(timestamp string) (string, string, error) {
	if len(timestamp) != 14 {
		return "", "", fmt.Errorf("length of timestamp is not as expected: %s", timestamp)
	}

	return timestamp[0:8], timestamp[8:14], nil
}

// StrToUnix converts time string of format "2016-01-11 06:12:33" to unix timestamp
//
// Params:
//   - timeStr: time string
//
// Returns:
//
//	(timestamp, err)
func StrToUnix(timeStr string) (int64, error) {
	t, err := xtime.Parse("2006-01-02 15:04:05 MST", timeStr+" CST")
	if err != nil {
		return 0, fmt.Errorf("wrong time string format")
	}
	return t.Unix(), nil
}

// UnixToStr converts unix timestamp to string of format "2006-01-02 15:04:05"
//
// Params:
//   - timestamp: unix timestamp
//
// Returns:
//
//	time string
func UnixToStr(timestamp int64) string {
	t := xtime.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// WaitTill waits until toTime
//
// Params:
//   - toTime: time to wait until. the number of seconds elapsed since January 1, 1970 UTC.
func WaitTill(toTime int64) {
	waitSecs := toTime - xtime.Now().Unix()
	if waitSecs > 0 {
		xtime.Sleep(xtime.Second * xtime.Duration(waitSecs))
	}
}

// CalcNextTime calculates the nearest time from now, given cycle and offset
//
// Params:
//   - cycle: cycle in seconds
//   - offset: offset of the next time; in seconds
//
// Return:
//   - timestamp of next time
func CalcNextTime(cycle int64, offset int64) int64 {
	current := xtime.Now().Unix()

	if current%cycle == 0 {
		return current + offset
	} else {
		return current - current%cycle + cycle + offset
	}
}
