package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"time"
)

const (
	timeFormatFull   = "2006-01-02 15:04:05.000"
	timeFormatNumber = "20060102150405"
)

func newZap(level string, w io.Writer) *zap.Logger {
	var encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     basicDateTimeEncoder,           // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(parseZapLevel(level))
	var writers []zapcore.WriteSyncer
	writers = append(writers, zapcore.AddSync(w))
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器配置
		zapcore.NewMultiWriteSyncer(writers...),
		atomicLevel, // 日志级别
	)
	var options []zap.Option
	options = append(options, zap.AddCaller())
	options = append(options, zap.AddCallerSkip(5))
	return zap.New(core, options...)
}

func parseZapLevel(str string) (level zapcore.Level) {
	str = strings.ToLower(str)
	switch str {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return
}

type DateTimeEncoder func(time.Time, zapcore.PrimitiveArrayEncoder)

func basicDateTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, timeFormatFull, enc)
}

func numberDateTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, timeFormatNumber, enc)
}

func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

func (e *DateTimeEncoder) UnmarshalText(text []byte) error {
	switch string(text) {
	case "numberdatetime":
		*e = numberDateTimeEncoder
	case "datetime":
		*e = basicDateTimeEncoder
	default:
		*e = basicDateTimeEncoder
	}
	return nil
}
