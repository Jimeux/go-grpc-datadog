package o11y

import (
	"context"
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var logger *zap.Logger

func InitLogger(name, env, ver, logPath string) func() {
	closeFile := func() {}
	logOut := []string{"stdout"}
	errOut := []string{"stderr"}
	if logPath != "" {
		filePath := path.Join(logPath, name+".log")
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			panic(err)
		}
		logOut = append(logOut, filePath)
		errOut = append(errOut, filePath)
		closeFile = func() { _ = file.Close() }
	}

	logger = zap.Must(zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "body",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "",
			CallerKey:      "caller",
			FunctionKey:    "function",
			StacktraceKey:  "error.stack",
			SkipLineEnding: false,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			// EncodeName:          nil,
			// NewReflectedEncoder: nil,
			// ConsoleSeparator:    "",
		},
		OutputPaths:      logOut,
		ErrorOutputPaths: errOut,
		InitialFields: map[string]any{
			"dd.service": name,
			"dd.env":     env,
			"dd.version": ver,
		},
	}.Build(zap.AddCallerSkip(1)))

	return func() {
		_ = logger.Sync()
		closeFile()
	}
}

func Err(ctx context.Context, err error, msg string /*, attrs ...zap.Field*/) {
	traceID, spanID := traceInfo(ctx)
	logger.Error(msg,
		zap.Error(err),
		zap.Uint64("dd.trace_id", traceID),
		zap.Uint64("dd.span_id", spanID),
		zap.String("error.msg", err.Error()),
		zap.String("error.type", fmt.Sprintf("%T", err)),
		/*zap.Object("attributes", zapcore.ObjectMarshalerFunc(func(inner zapcore.ObjectEncoder) error {
			for _, a := range attrs {
				a.AddTo(inner)
			}
			return nil
		})),*/
	)
}

func Info(ctx context.Context, msg string /*, attrs ...zap.Field*/) {
	traceID, spanID := traceInfo(ctx)
	logger.Info(msg,
		zap.Uint64("dd.trace_id", traceID),
		zap.Uint64("dd.span_id", spanID),
		/*zap.Object("attributes", zapcore.ObjectMarshalerFunc(func(inner zapcore.ObjectEncoder) error {
			for _, a := range attrs {
				a.AddTo(inner)
			}
			return nil
		})),*/
	)
}

func traceInfo(ctx context.Context) (uint64, uint64) {
	sc, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return 0, 0
	}
	return sc.Context().TraceID(), sc.Context().SpanID()
}
