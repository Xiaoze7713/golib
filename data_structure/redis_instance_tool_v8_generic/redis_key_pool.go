package redis_instance_tool

import (
	"context"
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/type_def"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"reflect"
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
	return m.selfKey + m.sep + fmt.Sprintf("%v", s)
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
	if m.nx {
		_, err = m.client.SetNX(ctx, m.SelfKey(key), value, ex).Result()
	} else {
		_, err = m.client.Set(ctx, m.SelfKey(key), value, ex).Result()
	}
	return
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

func (m *KeyPool[K, V]) GetAndUnmarshal(ctx context.Context, key K) (value *V, err error) {
	resBytes, err := m.client.Get(ctx, m.SelfKey(key)).Result()
	if err != nil {
		return nil, err
	}
	var tp V
	v := reflect.New(reflect.TypeOf(tp))
	value, ok := v.Interface().(*V)
	if !ok {
		xlog.Error("failed conv")
	}
	err = m.MarshalInterface.unmarshalFunction(resBytes, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
