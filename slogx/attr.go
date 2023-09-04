package sloghelper

import "log/slog"

var ErrorKey = "error"

func Error(err error) slog.Attr {
	return slog.String(ErrorKey, err.Error())
}

func ErrorValue(err error) slog.Value {
	return slog.StringValue(err.Error())
}
