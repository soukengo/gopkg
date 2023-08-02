package log

type Config struct {
	LoggerConfig `mapstructure:",squash"`
	Loggers      []LoggerConfig
}

type LoggerConfig struct {
	Level    string
	Name     string
	Split    bool
	Appender *AppenderConfig
}

type AppenderConfig struct {
	Console    bool
	FileName   string
	MaxSize    int  // 文件大小
	MaxAge     int  // 存储时长
	MaxBackups int  // 最大数量
	Compress   bool // 是否压缩
}

func Default() *Config {
	return &Config{
		LoggerConfig: LoggerConfig{
			Level: "info",
		},
	}
}
