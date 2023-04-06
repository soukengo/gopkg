package kafka

import (
	"context"
	"github.com/Shopify/sarama"
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
	consumer sarama.ConsumerGroup
	session  sarama.ConsumerGroupSession
	config   *ConsumerConfig
	handler  mq.Handler
	queue    chan *sarama.ConsumerMessage // 数据队列
	consumed sync.Once
	closed   sync.Once
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewConsumer(config *ConsumerConfig, handler mq.Handler, logger log.Logger) (*Consumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true
	if config.Workers <= 0 {
		config.Workers = defaultWorkers
	}
	consumer, err := sarama.NewConsumerGroup(config.Kafka.Brokers, config.Group, kafkaConfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &Consumer{
		logger:   logger,
		config:   config,
		consumer: consumer,
		handler:  handler,
		queue:    make(chan *sarama.ConsumerMessage, config.Workers),
		ctx:      ctx,
		cancel:   cancel,
	}
	return c, nil
}

// Start 消费
func (c *Consumer) Start() (err error) {
	c.consumed.Do(func() {
		go runtimes.TryCatch(func() {
			c.receive()
		})
		go runtimes.TryCatch(func() {
			c.dispatch()
		})
	})
	return
}

// dispatch 分发
func (c *Consumer) dispatch() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.queue:
			topic := msg.Topic
			// 移除配置的topic前缀
			topic = strings.TrimPrefix(topic, c.config.Kafka.TopicPrefix)
			err := c.handler.OnReceived(topic, msg.Value)
			if err != nil {
				log.Errorf("handler.OnReceived error: %v", err)
				continue
			}
			if c.session != nil {
				c.session.MarkMessage(msg, "")
			}
		}
	}
}

// receive 接收
func (c *Consumer) receive() {
	config := c.config
	topics := make([]string, len(config.Topics))
	for i := 0; i < len(config.Topics); i++ {
		topics[i] = config.Kafka.TopicPrefix + config.Topics[i]
	}
	for {
		if err := c.consumer.Consume(c.ctx, topics, c); err != nil {
			log.Errorf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if c.ctx.Err() != nil {
			return
		}
	}
}

func (c *Consumer) Close() error {
	c.closed.Do(func() {
		c.cancel()
		close(c.queue)
	})
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.session = session
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case err := <-c.consumer.Errors():
			c.logger.Helper().Errorf("consumer error(%v)", err)
		case message := <-claim.Messages():
			c.queue <- message
		case <-session.Context().Done():
			return nil
		}
	}
}
