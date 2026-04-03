package logging

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/SilentPlaces/basicauth.git/internal/config"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	"github.com/SilentPlaces/basicauth.git/internal/shared/observability"
	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	base zerolog.Logger
}

func NewZeroLogger(appCfg *config.AppConfig) appLogger.Logger {
	var output io.Writer = os.Stdout
	if strings.EqualFold(appCfg.LogFormat, "pretty") {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Logger().
		Level(parseLevel(appCfg.LogLevel))

	return &ZeroLogger{base: logger}
}

func (l *ZeroLogger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	event := withContext(l.base.Debug(), ctx)
	addFields(event, fields)
	event.Msg(message)
}

func (l *ZeroLogger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	event := withContext(l.base.Info(), ctx)
	addFields(event, fields)
	event.Msg(message)
}

func (l *ZeroLogger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	event := withContext(l.base.Warn(), ctx)
	addFields(event, fields)
	event.Msg(message)
}

func (l *ZeroLogger) Error(ctx context.Context, message string, err error, fields map[string]interface{}) {
	event := withContext(l.base.Error(), ctx)
	if err != nil {
		event = event.Err(err)
	}
	addFields(event, fields)
	event.Msg(message)
}

func addFields(event *zerolog.Event, fields map[string]interface{}) {
	for key, value := range fields {
		switch v := value.(type) {
		case string:
			event.Str(key, v)
		case int:
			event.Int(key, v)
		case int64:
			event.Int64(key, v)
		case bool:
			event.Bool(key, v)
		case time.Duration:
			event.Dur(key, v)
		default:
			event.Interface(key, v)
		}
	}
}

func withContext(event *zerolog.Event, ctx context.Context) *zerolog.Event {
	correlationID := observability.CorrelationIDFromContext(ctx)
	if correlationID != "" {
		event = event.Str("correlation_id", correlationID)
	}
	traceParent := observability.TraceParentFromContext(ctx)
	if traceParent != "" {
		event = event.Str("traceparent", traceParent)
	}
	return event
}

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
