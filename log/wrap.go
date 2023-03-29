package log

import (
	"github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/soukengo/gopkg/util/file"
	"path/filepath"
	"strings"
)

type leveledProvider struct {
	leveled map[Level]Provider
}

func newLeveledLogger(lv Level, cfg *LoggerConfig) (Provider, error) {
	leveled := make(map[Level]Provider)
	absPath, err := file.Abs(cfg.Appender.FileName)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(absPath)
	sourceName := filepath.Base(absPath)
	for _, level := range levelMap {
		if level < lv {
			continue
		}
		levelName := strings.ToLower(level.String())
		fileName := strings.Join([]string{dir, levelName, sourceName}, string(filepath.Separator))
		ac := &AppenderConfig{Pattern: cfg.Appender.Pattern, FileName: fileName, RotationCount: cfg.Appender.RotationCount}
		leveled[level] = newAdapter(level, ac)
	}
	return &leveledProvider{leveled: leveled}, nil
}

func (l *leveledProvider) Log(level log.Level, kvs ...interface{}) error {
	if impl, ok := l.leveled[Level(level)]; ok {
		return impl.Log(level, kvs...)
	}
	return nil
}

type singleProvider struct {
	impl log.Logger
}

func newSingleProvider(lv Level, cfg *LoggerConfig) Provider {
	var writer = newAppender(cfg.Appender)
	return &singleProvider{zap.NewLogger(newZap(lv.String(), writer))}
}

func (l *singleProvider) Log(level log.Level, kvs ...interface{}) error {
	return l.impl.Log(level, kvs...)
}

func newAdapter(level Level, appender *AppenderConfig) Provider {
	var writer = newAppender(appender)
	return zap.NewLogger(newZap(level.String(), writer))
}
