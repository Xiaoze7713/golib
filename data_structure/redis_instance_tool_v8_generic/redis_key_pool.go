package redis_instance_tool

import (
	"context"
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/v2/generic_conv"
	"git.singularity-ai.com/backend/library/v2/type_def"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"time"
)

type KeyPool[K type_def.BaseValueType, V any] struct {
	RedisToolBase
	nx     bool
	expire int // second
}

func NewKeyPool[K type_def.BaseValueType, V any](businessKey, sep string, client redis.Cmdable, nx bool, duration int) (keyPool *KeyPool[K, V], err error) {
	if businessKey == "" {
		xlog.Warn("key pool null business key")
	}
	if client == nil {
		err = errors.New("redis pool is nil")
	}
	keyPool = &KeyPool[K, V]{
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

func (m *KeyPool[K, V]) SelfKey(s K) (key string) {
	return m.Prefix() + fmt.Sprintf("%v", s)
}

func (m *KeyPool[K, V]) SelfSep() (sep string) {
	return m.sep
}

func (m *KeyPool[K, V]) Set(ctx context.Context, key K, valueIF V) (err error) {
	value, err := m.marshalFunction(valueIF)
	if err != nil {
		return err
	}
	//}
	ex := time.Duration(0)
	if m.expire != 0 {
		ex = time.Second * time.Duration(m.expire)
	}
	var res bool
	if m.nx {
		res, err = m.client.SetNX(ctx, m.SelfKey(key), value, ex).Result()
	} else {
		_, err = m.client.Set(ctx, m.SelfKey(key), value, ex).Result()
		res = true
	}
	if !res {
		return errors.New("set failed")
	}
	return
}

func (m *KeyPool[K, V]) Keys(ctx context.Context) (resList []string, err error) {
	if m.Prefix() == "" {
		return nil, WarnKeysCountUse
	}
	resList, err = m.client.Keys(ctx, m.Prefix()+"*").Result()
	if err != nil {
		return nil, err
	}
	return
}

func (m *KeyPool[K, V]) Count(ctx context.Context) (count int, err error) {
	if m.Prefix() == "" {
		return -1, WarnKeysCountUse
	}
	resList, err := m.client.Keys(ctx, m.Prefix()+"*").Result()
	if err != nil {
		return 0, err
	}
	return len(resList), nil
}

func (m *KeyPool[K, V]) Get(ctx context.Context, key K) (resStr string, err error) {
	resStr, err = m.client.Get(ctx, m.SelfKey(key)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return
}

func (m *KeyPool[K, V]) Del(ctx context.Context, key K) (err error) {
	_, err = m.client.Del(ctx, m.SelfKey(key)).Result()
	return
}

func (m *KeyPool[K, V]) GetAndUnmarshal(ctx context.Context, key K) (value V, err error) {
	resStr, err := m.client.Get(ctx, m.SelfKey(key)).Result()
	var v V
	if err != nil {
		return v, err
	}
	value, err1 := generic_conv.Conv[V](resStr, m.unmarshalFunction)
	if err1 != nil {
		xlog.Error(err1)
		return v, err1
	}
	return value, nil
}
