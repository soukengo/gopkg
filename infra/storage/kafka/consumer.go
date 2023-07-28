package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/dispatcher"
	"github.com/soukengo/gopkg/util/runtimes"
	"strings"

	"sync"
	"time"
)

type Consumer interface {
	Start() (err error)
	Stop() error
	Subscribe(topic string, handler HandlerFunc)
	Ack(m *Message)
}

type consumer struct {
	logger   log.Logger
	consumer sarama.ConsumerGroup
	session  sarama.ConsumerGroupSession
	cfg      *Config
	queue    chan *Message // 数据队列
	started  sync.Once
	stopped  sync.Once
	ctx      context.Context
	cancel   context.CancelFunc
	topics   map[string][]HandlerFunc
	mutex    sync.Mutex
}

func NewConsumer(cfg *Config, consumerCfg *ConsumerConfig, logger log.Logger) (Consumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true
	if consumerCfg.Workers <= 0 {
		consumerCfg.Workers = defaultWorkers
	}
	cg, err := sarama.NewConsumerGroup(cfg.Brokers, consumerCfg.GroupId, kafkaConfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &consumer{
		logger:   logger,
		cfg:      cfg,
		consumer: cg,
		queue:    make(chan *Message, consumerCfg.Workers),
		ctx:      ctx,
		cancel:   cancel,
	}
	return c, nil
}

func (c *consumer) Start() (err error) {
	c.started.Do(func() {
		runtimes.Async(c.receive)
		runtimes.Async(c.dispatch)
	})
	return
}

func (c *consumer) Stop() error {
	c.stopped.Do(func() {
		c.cancel()
		close(c.queue)
	})
	return nil
}

func (c *consumer) Subscribe(topic string, handler HandlerFunc) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	handlers := c.topics[topic]
	handlers = append(handlers, handler)
	c.topics[topic] = handlers
}

func (c *consumer) Ack(m *Message) {
	if c.session == nil {
		return
	}
	c.session.MarkOffset(m.topic, m.partition, m.offset+1, "")
}

// dispatch 分发
func (c *consumer) dispatch() {
	dispatcher.Dispatch(c.ctx, c.queue, func(msg *Message) {
		// 移除配置的topic前缀
		topic := strings.TrimPrefix(msg.originTopic, c.cfg.TopicPrefix)
		msg.topic = topic
		handlers, ok := c.topics[topic]
		if !ok {
			c.Ack(msg)
			return
		}
		for _, item := range handlers {
			var handler = item
			handler(msg)
		}
	})
}

// receive 接收数据
func (c *consumer) receive() {
	topics := make([]string, 0)
	for topic := range c.topics {
		topics = append(topics, c.cfg.TopicPrefix+topic)
	}
	for {
		if err := c.consumer.Consume(c.ctx, topics, c); err != nil {
			c.logger.Helper().Errorf("Error from consumer: %v", err)
			time.Sleep(time.Second)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if c.ctx.Err() != nil {
			return
		}
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.session = session
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case err := <-c.consumer.Errors():
			c.logger.Helper().Errorf("consumer error(%v)", err)
		case message := <-claim.Messages():
			c.queue <- &Message{key: message.Key, value: message.Value, originTopic: message.Topic, partition: message.Partition, offset: message.Offset}
		case <-session.Context().Done():
			return nil
		}
	}
}
