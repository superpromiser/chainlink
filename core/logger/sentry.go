package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/static"
)

const SentryFlushDeadline = 5 * time.Second

func init() {
	// If SENTRY_ENVIRONMENT is set, it will override everything. Otherwise infers from CHAINLINK_DEV.
	var sentryenv string
	if env := os.Getenv("SENTRY_ENVIRONMENT"); env != "" {
		sentryenv = env
	} else if os.Getenv("CHAINLINK_DEV") == "true" {
		sentryenv = "dev"
	} else {
		sentryenv = "prod"
	}
	// If SENTRY_DSN is set, it will override everything. Otherwise static.SentryDSN will be used.
	// If neither are set, sentry is disabled.
	var sentrydsn string
	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		sentrydsn = dsn
	} else {
		sentrydsn = static.SentryDSN
	}
	// If SENTRY_RELEASE is set, it will override everything. Otherwise, static.Version will be used.
	var sentryrelease string
	if release := os.Getenv("SENTRY_RELEASE"); release != "" {
		sentryrelease = release
	} else {
		sentryrelease = static.Version
	}
	err := sentry.Init(sentry.ClientOptions{
		// AttachStacktrace is needed to send stacktrace alongside panics
		AttachStacktrace: true,
		Dsn:              sentrydsn,
		Environment:      sentryenv,
		Release:          sentryrelease,
		// Enable printing of SDK debug messages.
		// Uncomment line below to debug sentry
		// Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

//TODO could include sentry event IDs with log lines: .With("sentryEvent",*EventID)
type sentryLogger struct {
	h Logger
}

func newSentryLogger(l Logger) Logger {
	return &sentryLogger{h: l.Helper(1)}
}

func (s *sentryLogger) With(args ...interface{}) Logger {
	return &sentryLogger{
		h: s.h.With(args...),
	}
}

func (s *sentryLogger) Named(name string) Logger {
	return &sentryLogger{
		h: s.h.Named(name),
	}
}

func (s *sentryLogger) NewRootLogger(lvl zapcore.Level) (Logger, error) {
	h, err := s.h.NewRootLogger(lvl)
	if err != nil {
		return nil, err
	}
	return &sentryLogger{
		h: h,
	}, nil
}

func (s *sentryLogger) SetLogLevel(level zapcore.Level) {
	s.h.SetLogLevel(level)
}

func (s *sentryLogger) Trace(args ...interface{}) {
	s.h.Trace(args...)
}

func (s *sentryLogger) Debug(args ...interface{}) {
	s.h.Debug(args...)
}

func (s *sentryLogger) Info(args ...interface{}) {
	s.h.Info(args...)
}

func (s *sentryLogger) Warn(args ...interface{}) {
	s.h.Warn(args...)
}

func (s *sentryLogger) Error(args ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelError)
	})
	hub.CaptureMessage(fmt.Sprintf("%v", args))
	s.h.Error(args...)
}

func (s *sentryLogger) Critical(args ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf("%v", args))
	s.h.Critical(args...)
}

func (s *sentryLogger) Panic(args ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf("%v", args))
	s.h.Panic(args...)
}

func (s *sentryLogger) Fatal(args ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf("%v", args))
	s.h.Fatal(args...)
}

func (s *sentryLogger) Tracef(format string, values ...interface{}) {
	s.h.Tracef(format, values...)
}

func (s *sentryLogger) Debugf(format string, values ...interface{}) {
	s.h.Debugf(format, values...)
}

func (s *sentryLogger) Infof(format string, values ...interface{}) {
	s.h.Infof(format, values...)
}

func (s *sentryLogger) Warnf(format string, values ...interface{}) {
	s.h.Warnf(format, values...)
}

func (s *sentryLogger) Errorf(format string, values ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelError)
	})
	hub.CaptureMessage(fmt.Sprintf(format, values...))
	s.h.Errorf(format, values...)
}

func (s *sentryLogger) Criticalf(format string, values ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf(format, values...))
	s.h.Criticalf(format, values...)
}

func (s *sentryLogger) Panicf(format string, values ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf(format, values...))
	s.h.Panicf(format, values...)
}

func (s *sentryLogger) Fatalf(format string, values ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf(format, values...))
	s.h.Fatalf(format, values...)
}

func (s *sentryLogger) Tracew(msg string, keysAndValues ...interface{}) {
	s.h.Tracew(msg, keysAndValues...)
}

func (s *sentryLogger) Debugw(msg string, keysAndValues ...interface{}) {
	s.h.Debugw(msg, keysAndValues...)
}

func (s *sentryLogger) Infow(msg string, keysAndValues ...interface{}) {
	s.h.Infow(msg, keysAndValues...)
}

func (s *sentryLogger) Warnw(msg string, keysAndValues ...interface{}) {
	s.h.Warnw(msg, keysAndValues...)
}

func (s *sentryLogger) Errorw(msg string, keysAndValues ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", toMap(keysAndValues))
		scope.SetLevel(sentry.LevelError)
	})
	hub.CaptureMessage(msg)
	s.h.Errorw(msg, keysAndValues...)
}

func (s *sentryLogger) CriticalW(msg string, keysAndValues ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", toMap(keysAndValues))
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(msg)
	s.h.CriticalW(msg, keysAndValues...)
}

func (s *sentryLogger) Panicw(msg string, keysAndValues ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", toMap(keysAndValues))
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(msg)
	s.h.Panicw(msg, keysAndValues...)
}

func (s *sentryLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", toMap(keysAndValues))
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(msg)
	s.h.Fatalw(msg, keysAndValues...)
}

func (s *sentryLogger) ErrorIf(err error, msg string) {
	if err != nil {
		sentry.CaptureException(err)
		s.h.Errorw(msg, "err", err)
	}
}

func (s *sentryLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		sentry.CaptureException(err)
		s.h.Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
	}
}

func (s *sentryLogger) Sync() error {
	return s.h.Sync()
}

func (s *sentryLogger) Helper(add int) Logger {
	return s.h.Helper(add)
}

func toMap(args ...interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); {
		// Make sure this element isn't a dangling key
		if i == len(args)-1 {
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string ignore it
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); ok {
			m[keyStr] = val
		}
		i += 2
	}
	return m
}
