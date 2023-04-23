package kafka

import (
	"context"
	kafka "github.com/Shopify/sarama"
	"github.com/google/uuid"
)

type Producer struct {
	cfg      *Config
	producer kafka.SyncProducer
}

func NewKafkaProducer(cfg *Config) (*Producer, error) {
	ins := &Producer{cfg: cfg}
	kc := kafka.NewConfig()
	kc.Producer.RequiredAcks = kafka.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                  // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := kafka.NewSyncProducer(cfg.Brokers, kc)
	if err != nil {
		return nil, err
	}
	ins.producer = pub
	return ins, nil
}

func (c *Producer) Close() error {
	return c.producer.Close()
}

func (c *Producer) Submit(ctx context.Context, topic string, values []Value) (err error) {
	var messages = make([]*kafka.ProducerMessage, len(values))
	for i, value := range values {
		messages[i] = &kafka.ProducerMessage{
			Key:   kafka.StringEncoder(uuid.New().String()),
			Topic: c.cfg.TopicPrefix + topic,
			Value: kafka.ByteEncoder(value),
		}
	}
	err = c.producer.SendMessages(messages)
	return
}
