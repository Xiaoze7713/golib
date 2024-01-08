package kafka_instance_tool

import (
	"context"
	"github.com/Xiaoze7713/golib/v3/data_structure/message_queue"
	"github.com/Xiaoze7713/golib/v3/kafka"
)

type KafkaProducer struct {
	kafka.Producer
	config *message_queue.MQConfig
}

func (m *KafkaProducer) Name() string {
	return m.config.Topic
}

func (m *KafkaProducer) Push(ctx context.Context, data []byte) error {
	m.ProduceChan <- data
	return nil
}

func (m *KafkaProducer) Pop(ctx context.Context, block bool) (data []byte, err error) {
	return nil, nil
}

func (m *KafkaProducer) Init(config *message_queue.MQConfig) (err error) {
	m.config = config
	return nil
}

func (m *KafkaProducer) Config() *message_queue.MQConfig {
	return m.config
}
