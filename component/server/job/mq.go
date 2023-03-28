package job

import (
	"context"
	"github.com/soukengo/gopkg/component/mq"
	"github.com/soukengo/gopkg/component/mq/kafka"
	"github.com/soukengo/gopkg/log"
	"sync"
)

type mqServer struct {
	consumer mq.Consumer
	cfg      *Config
	mqCfg    *mq.Config
	handlers map[string]Handler
	mutex    sync.Mutex
}

func NewMQServer(cfg *Config, mqCfg *mq.Config, logger log.Logger) (Server, error) {
	s := &mqServer{cfg: cfg, mqCfg: mqCfg, handlers: make(map[string]Handler)}
	consumer, err := kafka.NewConsumer(&kafka.ConsumerConfig{
		Kafka:  mqCfg.Kafka,
		Group:  cfg.GroupId,
		Topics: cfg.Topics,
	}, s, logger)
	if err != nil {
		return nil, err
	}
	s.consumer = consumer
	return s, nil
}

func (s *mqServer) Start(ctx context.Context) error {
	return s.consumer.Start()
}

func (s *mqServer) Stop(ctx context.Context) error {
	return s.consumer.Close()
}
func (s *mqServer) Register(topic string, h Handler) {
	s.handlers[topic] = h
}

func (s *mqServer) OnReceived(topic string, data []byte) error {
	h, ok := s.handlers[topic]
	if !ok {
		log.Warnf("OnReceived topic[%s] no handle found", topic)
	}
	ctx := context.Background()
	return h(ctx, topic, data)
}
