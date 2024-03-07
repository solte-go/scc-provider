package log

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

type Attr slog.Attr

type traceAttr struct {
	ID string `json:"id"`
}

func String(key, value string) Attr {
	return Attr(slog.String(key, value))
}

func Int64(key string, value int64) Attr {
	return Attr(slog.Int64(key, value))
}

func Int(key string, value int) Attr {
	return Attr(slog.Int(key, value))
}

func Uint64(key string, value uint64) Attr {
	return Attr(slog.Uint64(key, value))
}

func Float64(key string, value float64) Attr {
	return Attr(slog.Float64(key, value))
}

func Bool(key string, value bool) Attr {
	return Attr(slog.Bool(key, value))
}

func Time(key string, value time.Time) Attr {
	return Attr(slog.Time(key, value))
}

func Duration(key string, value time.Duration) Attr {
	return Attr(slog.Duration(key, value))
}

func Group(key string, args ...any) Attr {
	return Attr(slog.Group(key, args...))
}

func Any(key string, value any) Attr {
	return Attr(slog.Any(key, value))
}

func Code(code string) Attr {
	return String("code", code)
}

func Error(err error) Attr {
	return String("err", err.Error())
}

func Context(ctx any) Attr {
	buf, _ := json.Marshal(ctx) //nolint:errchkjson

	if len(buf) == 0 {
		return Any("context", nil)
	}

	return String("context", string(buf))
}

func TraceFromContext(ctx context.Context) Attr {
	if traceId, ok := ctx.Value(TraceKey).(string); ok {
		return Any("trace", traceAttr{ID: traceId})
	}

	return Any("trace", traceAttr{ID: ""})
}
