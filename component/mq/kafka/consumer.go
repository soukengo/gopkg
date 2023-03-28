package kafka

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/soukengo/gopkg/component/mq"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/runtimes"
	"strings"
	"sync"
)

const (
	defaultWorkers = 10
)

type Consumer struct {
	logger   log.Logger
	consumer *cluster.Consumer
	config   *ConsumerConfig
	handler  mq.Handler
	queue    chan *sarama.ConsumerMessage // 数据队列
	consumed sync.Once
	closed   sync.Once
	done     chan struct{}
}

func NewConsumer(config *ConsumerConfig, handler mq.Handler, logger log.Logger) (*Consumer, error) {
	kafkaConfig := cluster.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Group.Return.Notifications = true
	topics := make([]string, len(config.Topics))
	for i := 0; i < len(config.Topics); i++ {
		topics[i] = config.Kafka.TopicPrefix + config.Topics[i]
	}
	if config.Workers <= 0 {
		config.Workers = defaultWorkers
	}
	consumer, err := cluster.NewConsumer(config.Kafka.Brokers, config.Group, topics, kafkaConfig)
	if err != nil {
		return nil, err
	}
	c := &Consumer{
		logger:   logger,
		config:   config,
		consumer: consumer,
		handler:  handler,
		queue:    make(chan *sarama.ConsumerMessage, config.Workers),
		done:     make(chan struct{}),
	}
	return c, nil
}

// Start 消费
func (s *Consumer) Start() (err error) {
	s.consumed.Do(func() {
		go runtimes.TryCatch(func() {
			s.receive()
		})
		go runtimes.TryCatch(func() {
			s.dispatch()
		})
	})
	return
}

// dispatch 分发
func (s *Consumer) dispatch() {
	for {
		select {
		case <-s.done:
			return
		case msg := <-s.queue:
			topic := msg.Topic
			// 移除配置的topic前缀
			topic = strings.TrimPrefix(topic, s.config.Kafka.TopicPrefix)
			err := s.handler.OnReceived(topic, msg.Value)
			if err != nil {
				log.Errorf("handler.OnReceived error: %v", err)
				continue
			}
			s.consumer.MarkOffset(msg, "")
		}
	}
}

// receive 接收
func (s *Consumer) receive() {
	for {
		select {
		case <-s.done:
			return
		case err := <-s.consumer.Errors():
			s.logger.Helper().Errorf("consumer error(%v)", err)
		case n := <-s.consumer.Notifications():
			s.logger.Helper().Infof("consumer rebalanced(%v)", n)
		case msg := <-s.consumer.Messages():
			s.queue <- msg
		}
	}
}

func (s *Consumer) Close() error {
	s.closed.Do(func() {
		close(s.done)
		close(s.queue)
	})
	return nil
}
