package redis_instance_tool

import (
	"context"
	"errors"
	"github.com/golib/v2/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"time"
)

type Counter struct {
	RedisToolBase
	//RedisCli *redis.Client
	//selfKey  string
	//sep      string
}

func NewCounter(key, sep string, client redis.Cmdable) (counter *Counter, err error) {
	if key == "" {
		xlog.Warn("counter null business key")
		//err = errors.New("null business key")
	}
	if client == nil {
		err = errors.New("redis pool is nil")
	}
	counter = &Counter{
		RedisToolBase{
			client:  client,
			selfKey: key,
			sep:     sep,
		},
	}
	if counter.MarshalInterface == nil {
		defaultMarshal.Apply(&counter.RedisToolBase)
	}
	return
}

func (kp *Counter) SelfKey(s string) (key string) {
	return kp.Prefix() + s
}

func (kp *Counter) SelfSep() (sep string) {
	return kp.sep
}

func (m *Counter) Set(ctx context.Context, key string, value int64) (err error) {
	_, err = m.client.Set(ctx, m.SelfKey(key), value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (m *Counter) Add(ctx context.Context, key string, value int64) (newValue int64, err error) {
	res, err := m.client.IncrBy(ctx, m.SelfKey(key), value).Result()
	if err != nil {
		return -1, err
	}
	newValue = int64(res)
	return newValue, nil
}

func (m *Counter) Expire(ctx context.Context, key string, dur int64) (err error) {
	_, err = m.client.Expire(ctx, m.SelfKey(key), time.Second*time.Duration(dur)).Result()
	return nil
}

func (m *Counter) Del(ctx context.Context, key string) (err error) {
	_, err = m.client.Del(ctx, m.SelfKey(key)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (m *Counter) Reduce(ctx context.Context, key string, value int64) (newValue int64, err error) {
	newValue, err = m.client.DecrBy(ctx, m.SelfKey(key), value).Result()
	if err != nil {
		return -1, err
	}
	return newValue, nil
}

func (m *Counter) Get(ctx context.Context, key string) (value int64, err error) {
	value, err = m.client.IncrBy(ctx, m.SelfKey(key), 0).Result()
	if err != nil {
		return -1, err
	}
	return value, nil
}
