//go:generate mockery --inpackage --all --case underscore --dir .
package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type exitFunc func(int)

type Logger struct {
	exit exitFunc
	sl   *slog.Logger
}

type Level slog.Level

type contextKey string

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
	LevelFatal Level = 12

	TraceKey contextKey = "traceId"
)

func (l Level) String() string {
	str := func(base string, val Level) string {
		if val == 0 {
			return base
		}

		return fmt.Sprintf("%s%+d", base, val)
	}

	switch {
	case l < LevelInfo:
		return str("DEBUG", l-LevelDebug)
	case l < LevelWarn:
		return str("INFO", l-LevelInfo)
	case l < LevelError:
		return str("WARN", l-LevelWarn)
	case l < LevelFatal:
		return str("ERROR", l-LevelError)
	default:
		return str("FATAL", l-LevelFatal)
	}
}

func (l Level) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, l.String()), nil
}

func (l *Level) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	return l.parse(s)
}

func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Level) UnmarshalText(data []byte) error {
	return l.parse(string(data))
}

func (l *Level) parse(s string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("slog: level string %q: %w", s, err)
		}
	}()

	name := s
	offset := 0

	if i := strings.IndexAny(s, "+-"); i >= 0 {
		name = s[:i]

		offset, err = strconv.Atoi(s[i:])
		if err != nil {
			return err
		}
	}

	switch strings.ToUpper(name) {
	case "DEBUG":
		*l = LevelDebug
	case "INFO":
		*l = LevelInfo
	case "WARN":
		*l = LevelWarn
	case "ERROR":
		*l = LevelError
	case "FATAL":
		*l = LevelFatal
	default:
		return fmt.Errorf("unknown level %s", name) //nolint:goerr113
	}

	*l += Level(offset)

	return nil
}

func New() *Logger {
	level := LevelInfo
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		level = LevelInfo
	}

	return newLogger(level, os.Stderr)
}

func NewNop() *Logger {
	level := LevelInfo
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		level = LevelInfo
	}

	return newLogger(level, io.Discard)
}

func newLogger(level Level, w io.Writer) *Logger {
	return &Logger{
		exit: os.Exit,
		sl: slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:       slog.Level(level),
			ReplaceAttr: replaceAttr,
		})),
	}
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key != "level" {
		return a
	}

	l, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a
	}

	return slog.Any(a.Key, Level(l))
}

func (l *Logger) log(level Level, msg string, attrs ...Attr) {
	casted := make([]slog.Attr, len(attrs))
	for i := range attrs {
		casted[i] = slog.Attr(attrs[i])
	}

	l.sl.LogAttrs(context.Background(), slog.Level(level), msg, casted...)
}

func (l *Logger) Debug(msg string, attrs ...Attr) {
	l.log(LevelDebug, msg, attrs...)
}

func (l *Logger) Info(msg string, attrs ...Attr) {
	l.log(LevelInfo, msg, attrs...)
}

func (l *Logger) Warn(msg string, attrs ...Attr) {
	l.log(LevelWarn, msg, attrs...)
}

func (l *Logger) Error(msg string, attrs ...Attr) {
	l.log(LevelError, msg, attrs...)
}

func (l *Logger) Fatal(msg string, attrs ...Attr) {
	l.log(LevelFatal, msg, attrs...)
	l.exit(1)
}
