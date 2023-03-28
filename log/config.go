package log

var (
	Default = &Config{
		LoggerConfig: LoggerConfig{
			Level: "info",
		},
	}
)

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
	Console       bool
	FileName      string
	Pattern       string
	RotationCount uint
}
