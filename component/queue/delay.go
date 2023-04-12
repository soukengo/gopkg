package queue

import (
	"time"
)

const (
	delayIdleTime      = time.Millisecond * 100
	delayConsumerBatch = 20
	defaultWorkers     = 1
)

type HandlerFunc func(topic Topic, value string)

type Queue interface {
	// Start 启动队列
	Start() error
	// Stop 停止队列
	Stop() error
	// Subscribe 订阅TOPIC
	Subscribe(topic Topic, handler HandlerFunc)
	// Submit 往队列提交数据
	Submit(topic Topic, values []string) error
}

type Delayed interface {
	Queue
	// SubmitDelay 提交延时数据
	SubmitDelay(topic Topic, values []string, delayMills int64, overwritten bool) error
	// RemoveDelay 从延时队列移除任务
	RemoveDelay(topic Topic, values []string) (deleted bool, err error)
}
