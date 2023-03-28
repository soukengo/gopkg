package http

import (
	"time"
)

var (
	client = NewClient(WithTimeout(time.Second * 10))
)

func Default() *Client {
	return client
}
