package runtimes

import (
	"fmt"
	"runtime"
	"time"
)

func Stack(err any) string {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	return fmt.Sprintf("\n%s\n%s\n%s", time.Now().Format("2006-01-02 15:04:05"), err, buf)
}
