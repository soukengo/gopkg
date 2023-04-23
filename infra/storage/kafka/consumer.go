package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/runtimes"

	//"github.com/soukengo/gopkg/util/runtimes"
	"strings"
	"sync"
	"time"
)

type Consumer struct {
	logger      log.Logger
	consumer    sarama.ConsumerGroup
	session     sarama.ConsumerGroupSession
	cfg         *Config
	consumerCfg *ConsumerConfig
	queue       chan *sarama.ConsumerMessage // 数据队列
	started     sync.Once
	stopped     sync.Once
	ctx         context.Context
	cancel      context.CancelFunc
	handlers    map[string]HandlerFunc
	mutex       sync.Mutex
}

func NewKafkaConsumer(cfg *Config, consumerCfg *ConsumerConfig, logger log.Logger) (*Consumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true
	if consumerCfg.Workers <= 0 {
		consumerCfg.Workers = defaultWorkers
	}
	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, consumerCfg.GroupId, kafkaConfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &Consumer{
		logger:      logger,
		cfg:         cfg,
		consumerCfg: consumerCfg,
		consumer:    consumer,
		queue:       make(chan *sarama.ConsumerMessage, consumerCfg.Workers),
		ctx:         ctx,
		cancel:      cancel,
		handlers:    make(map[string]HandlerFunc),
	}
	return c, nil
}

func (c *Consumer) Start() (err error) {
	c.started.Do(func() {
		go runtimes.TryCatch(func() {
			c.receive()
		})
		go runtimes.TryCatch(func() {
			c.dispatch()
		})
	})
	return
}

func (c *Consumer) Stop() error {
	c.stopped.Do(func() {
		c.cancel()
		close(c.queue)
	})
	return nil
}

func (c *Consumer) Subscribe(topic string, handler HandlerFunc) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers[topic] = handler
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
			topic = strings.TrimPrefix(topic, c.cfg.TopicPrefix)
			handler, ok := c.handlers[topic]
			if ok {
				err := handler(topic, msg.Value)
				if err != nil {
					log.Errorf("handler.OnReceived error: %v", err)
					continue
				}
			}
			if c.session != nil {
				c.session.MarkMessage(msg, "")
			}
		}
	}
}

// receive 接收
func (c *Consumer) receive() {
	config := c.consumerCfg
	topics := make([]string, len(config.Topics))
	for i := 0; i < len(config.Topics); i++ {
		topics[i] = c.cfg.TopicPrefix + string(config.Topics[i])
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
