package redis_instance_tool

import (
	"errors"
	"github.com/golib/v3/data_structure/message_queue"
	"github.com/golib/v3/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
)

type RedisMQ[T any] struct {
	*RedisQueue[T]
	config *message_queue.MQConfig
}

func (m *RedisMQ[T]) Config() *message_queue.MQConfig {
	return m.config
}

func NewRedisMQ[T any](msConfig *message_queue.MQConfig, client redis.Cmdable) (queue message_queue.MQInstance[T], err error) {
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
	redisQueue, err := NewQueue[T](msConfig.Topic, client)
	redisMessageQueue := &RedisMQ[T]{
		RedisQueue: redisQueue,
		config:     msConfig,
	}
	return redisMessageQueue, nil
}
