package redis_instance_tool

import (
	"errors"
	"git.singularity-ai.com/backend/library/data_structure/message_queue"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
)

type RedisMQ struct {
	*RedisQueue
	config *message_queue.MQConfig
}

func (m *RedisMQ) Config() *message_queue.MQConfig {
	return m.config
}

func NewRedisMQ(msConfig *message_queue.MQConfig, client redis.Cmdable) (queue message_queue.MQInstance, err error) {
	if msConfig.Topic == "" {
		xlog.Warn("queue null business key")
		//err = errors.New("null business key")
		//xlog.Error(err)
		return
	}
	if client == nil {
		err = errors.New("redis queue is nil")
		xlog.Error(err)
		return
	}
	redisQueue, err := NewQueue(msConfig.Topic, client)
	redisMessageQueue := &RedisMQ{
		RedisQueue: redisQueue,
		config:     msConfig,
	}
	return redisMessageQueue, nil
}
