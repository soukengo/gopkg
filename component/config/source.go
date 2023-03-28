package config

import "io"

type Source interface {
	From() (io.Reader, error)
	ConfigType() string
	Close() error
}
