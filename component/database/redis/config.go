package redis

import (
	"time"
)

const (
	Max = "+inf"
	Min = "-inf"
)

type Config struct {
	Addr         []string
	Password     string
	Prefix       string
	Active       int
	Idle         int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}
