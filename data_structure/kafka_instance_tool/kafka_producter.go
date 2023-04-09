package kafka_instance_tool

import (
	"context"
	"git.singularity-ai.com/backend/library/v2/data_structure/message_queue"
	"git.singularity-ai.com/backend/library/v2/kafka"
	json "github.com/json-iterator/go"
)

type KafkaProducer struct {
	kafka.Producer
	config *message_queue.MQConfig
}

func (m *KafkaProducer) Name() string {
	return m.config.Topic
}

func (m *KafkaProducer) Push(ctx context.Context, i interface{}) error {
	msg, err := json.MarshalToString(i)
	if err != nil {
		return err
	}
	m.ProduceChan <- []byte(msg)
	return nil
}

func (m *KafkaProducer) Pop(ctx context.Context, block bool) (i interface{}, err error) {
	return nil, nil
}

func (m *KafkaProducer) Init(config *message_queue.MQConfig) (err error) {
	m.config = config
	return nil
}

func (m *KafkaProducer) Config() *message_queue.MQConfig {
	return m.config
}
