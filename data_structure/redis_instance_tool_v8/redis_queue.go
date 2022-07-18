package redis_instance_tool

import (
	"context"
	"errors"
	"git.singularity-ai.com/backend/library/data_structure/message_queue"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/xutils/lib-common/utils"
)

type RedisQueue struct {
	RedisToolBase
}

func (m *RedisQueue) Init(msConfig *message_queue.MQConfig) error {
	m.selfKey = msConfig.Topic
	m.config = msConfig
	return nil
}

func (m *RedisQueue) Config() *message_queue.MQConfig {
	return m.config
}

func NewQueue(msConfig *message_queue.MQConfig, client *redis.Client) (queue message_queue.MQInstance, err error) {
	if msConfig.Topic == "" {
		xlog.Warn("queue null business key")
		//err = errors.New("null business key")
		//xlog.Error(err)
		return
	}
	if client == nil {
		err = errors.New("redis pool is nil")
		xlog.Error(err)
		return
	}
	queue = &RedisQueue{
		RedisToolBase{
			client: client,
		},
	}
	err = queue.Init(msConfig)
	return
}

func (m *RedisQueue) Name() string {
	return m.selfKey
}

func (m *RedisQueue) QueueName() string {
	return m.selfKey
}

func (m *RedisQueue) Push(ctx context.Context, valueIf interface{}) (err error) {
	value, err := jsoniter.MarshalToString(valueIf)
	_, err = m.client.RPush(ctx, m.QueueName(), value).Result()
	//xlog.Debugf("push queue %v,  count %v  value %v", m.QueueName(), r, value)
	return
}

func (m *RedisQueue) Pop(ctx context.Context, block bool) (interface{}, error) {
	res, err := m.pop(ctx, block)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *RedisQueue) pop(ctx context.Context, block bool) (res string, err error) {
	resList := []string{}
	if block {
		resList, err = m.client.BLPop(ctx, 0, m.QueueName()).Result()
	} else {
		res, err = m.client.LPop(ctx, m.QueueName()).Result()
		resList = append(resList, m.QueueName(), res)
	}
	xlog.Debug(utils.MustString(resList))
	if err != nil {
		return "", err
	}
	if len(resList) == 2 {
		return resList[1], nil
	}
	return "", nil
}
