package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestLog(t *testing.T) {
	t.Parallel()

	const msg = "msg"

	type record struct {
		Level   Level     `json:"level"`
		Msg     string    `json:"msg"`
		Code    string    `json:"code"`
		Err     string    `json:"err"`
		Context string    `json:"context"`
		Trace   traceAttr `json:"trace"`
	}

	cases := []struct {
		desc          string
		level         Level
		msg           string
		code          string
		err           error
		context       any
		tracedContext context.Context
		record        record
	}{
		{
			desc:  "debug with code",
			level: LevelDebug,
			msg:   msg,
			code:  "code",
			record: record{
				Level: LevelDebug,
				Msg:   msg,
				Code:  "code",
			},
		},
		{
			desc:  "info with context",
			level: LevelInfo,
			msg:   msg,
			context: struct {
				Prop string
			}{
				Prop: "prop",
			},
			record: record{
				Level:   LevelInfo,
				Msg:     msg,
				Context: `{"Prop":"prop"}`,
			},
		},
		{
			desc:  "warn with error",
			level: LevelWarn,
			msg:   msg,
			err:   io.EOF,
			record: record{
				Level: LevelWarn,
				Msg:   msg,
				Err:   io.EOF.Error(),
			},
		},
		{
			desc:  "error with error and code",
			level: LevelError,
			msg:   msg,
			code:  "code",
			err:   io.EOF,
			record: record{
				Level: LevelError,
				Msg:   msg,
				Code:  "code",
				Err:   io.EOF.Error(),
			},
		},
		{
			desc:  "fatal",
			level: LevelFatal,
			msg:   msg,
			record: record{
				Level: LevelFatal,
				Msg:   msg,
			},
		},
		{
			desc:          "info with traceId in context",
			level:         LevelInfo,
			msg:           msg,
			tracedContext: context.WithValue(context.Background(), TraceKey, "trace"),
			record: record{
				Level: LevelInfo,
				Msg:   msg,
				Trace: traceAttr{
					ID: "trace",
				},
			},
		},
	}

	for i := range cases {
		c := cases[i]

		t.Run(c.desc, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			exit := newMockExitFunc(t)
			if c.level == LevelFatal {
				exit.On("Execute", 1).Once()
			}

			l := newLogger(LevelDebug, &buf)
			l.exit = exit.Execute

			var fn func(string, ...Attr)

			switch c.level {
			case LevelDebug:
				fn = l.Debug
			case LevelInfo:
				fn = l.Info
			case LevelWarn:
				fn = l.Warn
			case LevelError:
				fn = l.Error
			case LevelFatal:
				fn = l.Fatal
			default:
				t.Logf("unknown level %s", c.level)
				t.FailNow()
			}

			var opts []Attr

			if c.code != "" {
				opts = append(opts, Code(c.code))
			}

			if c.err != nil {
				opts = append(opts, Error(c.err))
			}

			if c.context != nil {
				opts = append(opts, Context(c.context))
			}

			if c.tracedContext != nil {
				opts = append(opts, TraceFromContext(c.tracedContext))
			}

			fn(c.msg, opts...)

			var actual record
			if err := json.Unmarshal(buf.Bytes(), &actual); err != nil {
				t.Log(err)
				t.FailNow()
			}

			require.Equal(t, c.record, actual)

			if c.level == LevelFatal {
				exit.AssertCalled(t, "Execute", 1)
			}
		})
	}
}
