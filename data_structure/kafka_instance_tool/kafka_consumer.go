package kafka_instance_tool

import (
	"context"
	"github.com/Xiaoze7713/golib/v3/data_structure/kafka"
	"github.com/Xiaoze7713/golib/v3/data_structure/message_queue"
)

type KafkaConsumer struct {
	kafka.Consumer
	config *message_queue.MQConfig
}

func (m *KafkaConsumer) Name() string {
	return m.config.Topic
}

func (m *KafkaConsumer) Push(ctx context.Context, data []byte) error {
	return nil
}

func (m *KafkaConsumer) Pop(ctx context.Context, block bool) (data []byte, err error) {
	kafkaData := <-m.ConsumeChan
	return kafkaData.Value, nil
}

func (m *KafkaConsumer) Init(config *message_queue.MQConfig) (err error) {
	m.config = config
	return nil
}

func (m *KafkaConsumer) Config() *message_queue.MQConfig {
	return m.config
}
