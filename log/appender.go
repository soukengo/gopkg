package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"os"
	"strings"
	"time"
)

const (
	defaultLogFilePattern = "-%Y-%m-%d"
)

type appender struct {
	*rotatelogs.RotateLogs
}

func newAppender(c *AppenderConfig) io.Writer {
	if c == nil || c.Console {
		return os.Stdout
	}
	if len(c.Pattern) == 0 {
		c.Pattern = defaultLogFilePattern
	}
	rotateLogs, err := rotatelogs.New(
		c.FileName+c.Pattern,
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(c.FileName),
		// WithRotationTime设置日志分割的时间，这里设置为一天分割一次
		parseRotationTimeFromPattern(c.Pattern),
		// 设置文件最大分割数，超过此数量则删除多余文件。默认为0,不删除
		rotatelogs.WithRotationCount(c.RotationCount),
	)
	if err != nil {
		panic(err)
	}
	return &appender{
		RotateLogs: rotateLogs,
	}
}

// parseRotationTimeFromPattern 根据日志规则解析得到分割选项
func parseRotationTimeFromPattern(pattern string) rotatelogs.Option {
	// 默认按天分割
	duration := time.Hour * 24
	if strings.HasSuffix(pattern, "%H") { // 按小时分割
		duration = time.Hour
	} else if strings.HasSuffix(pattern, "%M") { // 按分钟分割
		duration = time.Minute
	}
	return rotatelogs.WithRotationTime(duration)
}
