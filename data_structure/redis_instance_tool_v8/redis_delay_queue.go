package redis_instance_tool

import (
	"context"
	"fmt"
	"git.singularity-ai.com/backend/library/data_structure/message_queue"
	"git.singularity-ai.com/backend/library/utils"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
)

type RedisDelayQueue struct {
	RedisToolBase
	limit int64
}

func (m *RedisDelayQueue) Init(msConfig *message_queue.MQConfig) error {
	m.selfKey = msConfig.Topic
	m.limit = 100
	m.config = msConfig
	if m.MarshalInterface == nil {
		defaultMarshal.Apply(&m.RedisToolBase)
	}
	return nil
}

func (m *RedisDelayQueue) Config() *message_queue.MQConfig {
	return m.config
}

func NewDelayQueue(client redis.Cmdable) (dQueue *RedisDelayQueue, err error) {
	dQueue = &RedisDelayQueue{
		RedisToolBase: RedisToolBase{
			client: client,
		},
	}
	return
}

func (m *RedisDelayQueue) Name() string {
	return m.selfKey
}

func (m *RedisDelayQueue) QueueName() string {
	return m.selfKey
}

func (m *RedisDelayQueue) ToKey(key string) string {
	return m.selfKey + "_kv_" + key
}

// json 写入， json解析

func (m *RedisDelayQueue) Add(ctx context.Context, key string, valueIf interface{}, delayDur int64, expireDur int64) (err error) {
	curTimeMs := time.Now().UnixMilli()
	value, err := jsoniter.MarshalToString(valueIf)
	if expireDur < delayDur {
		expireDur = delayDur + 1
	}
	timeMs := curTimeMs + utils.SecondToMs(delayDur)
	_, err = m.client.Set(ctx, m.ToKey(key), value, time.Second*time.Duration(expireDur)).Result()
	if err != nil {
		xlog.Errorf("set err %v %v %v", key, value, err)
		return err
	}
	_, err = m.client.ZAdd(ctx, m.QueueName(), &redis.Z{
		Score:  float64(timeMs),
		Member: m.ToKey(key),
	}).Result()
	if err != nil {
		xlog.Errorf("zadd err %v %v %v", key, value, err)
		return err
	}
	return
}

func (m *RedisDelayQueue) Count(ctx context.Context, timeMsStart, timeMsEnd int64) (count int, err error) {
	min := fmt.Sprintf("%d", timeMsStart)
	max := fmt.Sprintf("%d", timeMsEnd)
	res, err := m.client.ZCount(ctx, m.QueueName(), min, max).Result()
	if err != nil {
		return 0, err
	}
	count, err = strconv.Atoi(utils.UnMarshalToString(res))
	if err != nil {
		return 0, err
	}
	return
}

func (m *RedisDelayQueue) PopL(ctx context.Context, timeMsStart, timeMsEnd int64) (resList []string, err error) {
	queueKeyList, err := m.client.ZRangeByScore(ctx, m.QueueName(), &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", timeMsStart),
		Max:    fmt.Sprintf("%d", timeMsEnd),
		Offset: 0,
		Count:  m.limit,
	}).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(queueKeyList) == 0 {
		return nil, nil
	}
	resList = []string{}
	var invalidKeys []interface{}
	for _, key := range queueKeyList {
		value, err1 := m.client.Get(ctx, key).Result()
		if err1 != nil {
			if err1 == redis.Nil {
				xlog.Debugf("invalid key add %v", key)
				invalidKeys = append(invalidKeys, key)
			}
			xlog.Errorf("%v", err1)
			continue
		}
		_, err1 = m.client.Del(ctx, key).Result()
		if err1 != nil {
			xlog.Errorf("%v", err1)
			continue
		}
		resList = append(resList, value)
		invalidKeys = append(invalidKeys, key)
	}
	xlog.Debugf("res_list=%v", len(resList))
	_, err = m.client.ZRem(ctx, m.QueueName(), invalidKeys...).Result()
	if err != nil {
		xlog.Error(invalidKeys...)
		return nil, err
	}
	return resList, nil
}
