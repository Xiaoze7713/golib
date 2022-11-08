package redis_instance_tool

import (
	"context"
	"errors"
	"git.singularity-ai.com/backend/library/type_def"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"reflect"
)

type RedisQueue[T type_def.PtrValueType] struct {
	*RedisToolBase
}

func (m *RedisQueue[T]) Init(key string) error {
	m.selfKey = key
	return nil
}

func NewQueue[T type_def.PtrValueType](queueKey string, client redis.Cmdable) (queue *RedisQueue[T], err error) {
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

func (m *RedisQueue[T]) Pop(ctx context.Context, block bool) (T, error) {
	t, err := m.pop(ctx, block)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (m *RedisQueue[T]) pop(ctx context.Context, block bool) (t2 T, err error) {
	resList := []string{}
	if block {
		resList, err = m.client.BLPop(ctx, 0, m.QueueName()).Result()
	} else {
		var res string
		res, err = m.client.LPop(ctx, m.QueueName()).Result()
		resList = append(resList, m.QueueName(), res)
	}
	if err != nil {
		return nil, err
	}
	if len(resList) == 2 {
		var t T
		v := reflect.New(reflect.TypeOf(t))
		t2, ok := v.Interface().(T)
		if !ok {
			xlog.Error("failed conv")
		}
		err = m.RedisToolBase.unmarshalFunction(resList[1], t2)
		if err != nil {
			xlog.Error(err)
		}
		return t2, nil
	}
	return nil, nil
}
