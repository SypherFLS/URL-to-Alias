package sl

import "golang.org/x/exp/slog"


func Err(err error) slog.Attr {
	return slog.Attr{
		Key: "Error",
		Value: slog.StringValue(err.Error()),
	}
}