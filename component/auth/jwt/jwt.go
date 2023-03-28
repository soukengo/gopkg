package jwt

import "time"

type Config struct {
	AppKey  string
	Expires time.Duration
}
