package redis_instance_tool

import (
	"context"
	"errors"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"time"
)

type KeyPool struct {
	RedisToolBase
	nx     bool
	expire int // second
}

func NewKeyPool(businessKey, sep string, client redis.Cmdable, nx bool, duration int) (keyPool *KeyPool, err error) {
	if businessKey == "" {
		xlog.Warn("key pool null business key")
	}
	if client == nil {
		err = errors.New("redis pool is nil")
	}
	keyPool = &KeyPool{
		RedisToolBase{
			client:  client,
			selfKey: businessKey,
			sep:     sep,
		},
		nx,
		duration,
	}
	if keyPool.MarshalInterface == nil {
		defaultMarshal.Apply(&keyPool.RedisToolBase)
	}
	return
}

func (m *KeyPool) SelfKey(s string) (key string) {
	return m.selfKey + m.sep + s
}

func (m *KeyPool) SelfSep() (sep string) {
	return m.sep
}

func (m *KeyPool) Set(ctx context.Context, key string, valueIF interface{}) (err error) {
	value, err := m.MarshalInterface.marshalFunction(valueIF)
	if err != nil {
		return err
	}
	//}
	ex := time.Duration(0)
	if m.expire != 0 {
		ex = time.Second * time.Duration(m.expire)
	}
	if m.nx {
		_, err = m.client.SetNX(ctx, m.SelfKey(key), value, ex).Result()
	} else {
		_, err = m.client.Set(ctx, m.SelfKey(key), value, ex).Result()
	}
	return
}

func (m *KeyPool) Get(ctx context.Context, key string) (resStr string, err error) {
	resStr, err = m.client.Get(ctx, m.SelfKey(key)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return
}

func (m *KeyPool) Del(ctx context.Context, key string) (err error) {
	_, err = m.client.Del(ctx, m.SelfKey(key)).Result()
	return
}

func (m *KeyPool) GetAndUnmarshal(ctx context.Context, key string, ifc interface{}) (err error) {
	resBytes, err := m.client.Get(ctx, m.SelfKey(key)).Result()
	if err != nil {
		return err
	}
	err = m.MarshalInterface.unmarshalFunction(resBytes, ifc)
	if err != nil {
		return err
	}
	return nil
}
