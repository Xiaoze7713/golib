package redis_instance_tool

import (
	"context"
	"errors"
	"github.com/golib/v3/generic_conv"
	"github.com/golib/v3/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
)

type RedisQueue[T any] struct {
	*RedisToolBase
}

func (m *RedisQueue[T]) Init(key string) error {
	m.selfKey = key
	return nil
}

func NewQueue[T any](queueKey string, client redis.Cmdable) (queue *RedisQueue[T], err error) {
	if queueKey == "" {
		xlog.Warn("queue null business key")
		return
	}
	if client == nil {
		err = errors.New("redis pool is nil")
		xlog.Error(err)
		return
	}
	queue = &RedisQueue[T]{
		&RedisToolBase{
			client: client,
		},
	}
	if queue.MarshalInterface == nil {
		defaultMarshal.Apply(queue.RedisToolBase)
	}
	err = queue.Init(queueKey)
	return
}

func (m *RedisQueue[T]) Name() string {
	return m.selfKey
}

func (m *RedisQueue[T]) QueueName() string {
	return m.selfKey
}

func (m *RedisQueue[T]) Push(ctx context.Context, valuePtr T) (err error) {
	value, err := m.marshalFunction(valuePtr)
	_, err = m.client.RPush(ctx, m.QueueName(), value).Result()
	//xlog.Debugf("push queue %v,  count %v  value %v", m.QueueName(), r, value)
	return
}

func (m *RedisQueue[T]) Len(ctx context.Context) (length int64, err error) {
	length, err = m.client.LLen(ctx, m.QueueName()).Result()
	return
}

func (m *RedisQueue[T]) Pop(ctx context.Context, block bool) (T, error) {
	t, err := m.pop(ctx, block)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (m *RedisQueue[T]) pop(ctx context.Context, block bool) (res T, err error) {
	resList := []string{}
	if block {
		resList, err = m.client.BLPop(ctx, 0, m.QueueName()).Result()
	} else {
		var res string
		res, err = m.client.LPop(ctx, m.QueueName()).Result()
		resList = append(resList, m.QueueName(), res)
	}
	var t T
	if err != nil {
		return t, err
	}
	if len(resList) == 2 {
		res, err := generic_conv.Conv[T](resList[1], m.unmarshalFunction)
		if err != nil {
			xlog.Error(err)
		}
		return res, nil
	}
	return t, nil
}
