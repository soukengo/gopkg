package kafka

import (
	"context"
	kafka "github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/soukengo/gopkg/component/mq"
)

type Producer struct {
	brokers     []string
	topicPrefix string
	producer    kafka.SyncProducer
}

func NewProducer(c *mq.Kafka) (*Producer, error) {
	ins := &Producer{brokers: c.Brokers, topicPrefix: c.TopicPrefix}
	kc := kafka.NewConfig()
	kc.Producer.RequiredAcks = kafka.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                  // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := kafka.NewSyncProducer(ins.brokers, kc)
	if err != nil {
		return nil, err
	}
	ins.producer = pub
	return ins, nil
}

func (c *Producer) Close() error {
	return c.producer.Close()
}

func (c *Producer) Send(ctx context.Context, topic string, bytes []byte) (err error) {
	m := &kafka.ProducerMessage{
		Key:   kafka.StringEncoder(uuid.New().String()),
		Topic: c.topicPrefix + topic,
		Value: kafka.ByteEncoder(bytes),
	}
	_, _, err = c.producer.SendMessage(m)
	return
}
