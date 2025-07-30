package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type contextKey string

const TraceIdKey contextKey = "trace_id"

var Log *zerolog.Logger

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	IsDev      string
}

func InitLogger(config LoggerConfig) {
	Log = NewLogger(config)
}

func NewLogger(config LoggerConfig) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	var writer io.Writer

	if config.IsDev == "development" {
		if strings.Contains(config.Filename, "app.log") {
			writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		} else {
			writer = PrettyJSONWriter{Writer: os.Stdout}
		}
	} else {
		writer = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}

type PrettyJSONWriter struct {
	Writer io.Writer
}

func (w PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, p, "", "  ")
	if err != nil {
		return w.Writer.Write(p)
	}

	return w.Writer.Write(prettyJSON.Bytes())
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIdKey).(string); ok {
		return traceID
	}

	return ""
}
