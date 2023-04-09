package kafka_instance_tool

import (
	"context"
	"git.singularity-ai.com/backend/library/v2/data_structure/message_queue"
	"git.singularity-ai.com/backend/library/v2/kafka"
)

type KafkaConsumer struct {
	kafka.Consumer
	config *message_queue.MQConfig
}

func (m *KafkaConsumer) Name() string {
	return m.config.Topic
}

func (m *KafkaConsumer) Push(ctx context.Context, i interface{}) error {
	return nil
}

func (m *KafkaConsumer) Pop(ctx context.Context, block bool) (i interface{}, err error) {
	kafkaData := <-m.ConsumeChan
	return string(kafkaData.Value), nil
}

func (m *KafkaConsumer) Init(config *message_queue.MQConfig) (err error) {
	m.config = config
	return nil
}

func (m *KafkaConsumer) Config() *message_queue.MQConfig {
	return m.config
}
