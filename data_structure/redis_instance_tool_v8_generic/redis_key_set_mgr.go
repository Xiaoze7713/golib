package redis_instance_tool

import (
	"context"
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/v2/generic_conv"
	"git.singularity-ai.com/backend/library/v2/type_def"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
)

type KeySet[K type_def.BaseValueType, V any] struct {
	RedisToolBase
	nx     bool
	expire int // second
}

func NewKeySet[K type_def.BaseValueType, V any](businessKey, sep string, client redis.Cmdable) (keyPool *KeySet[K, V], err error) {
	if businessKey == "" {
		xlog.Warn("key pool null business key")
	}
	if client == nil {
		err = errors.New("redis pool is nil")
	}
	keyPool = &KeySet[K, V]{
		RedisToolBase: RedisToolBase{
			client:  client,
			selfKey: businessKey,
			sep:     sep,
		},
		//nx,
		//duration,
	}
	if keyPool.MarshalInterface == nil {
		defaultMarshal.Apply(&keyPool.RedisToolBase)
	}
	return
}

func (m *KeySet[K, V]) SelfKey(s K) (key string) {
	return m.Prefix() + fmt.Sprintf("%v", s)
}

func (m *KeySet[K, V]) SelfSep() (sep string) {
	return m.sep
}

func (m *KeySet[K, V]) Add(ctx context.Context, key K, valueIFs ...V) (err error) {
	valueList := []interface{}{}
	for _, vIf := range valueIFs {
		value, err1 := m.marshalFunction(vIf)
		if err1 != nil {
			return err1
		}
		valueList = append(valueList, value)
	}
	_, err = m.client.SAdd(ctx, m.SelfKey(key), valueList...).Result()
	return err
}

func (m *KeySet[K, V]) Count(ctx context.Context, key K) (count int, err error) {
	cnt, err := m.client.SCard(ctx, m.SelfKey(key)).Result()
	if err != nil {
		return 0, err
	}
	return int(cnt), nil
}

func (m *KeySet[K, V]) ExistMember(ctx context.Context, key K, valueIF V) (ok bool, err error) {
	value, err := m.marshalFunction(valueIF)
	if err != nil {
		return false, err
	}
	ok, err = m.client.SIsMember(ctx, m.SelfKey(key), value).Result()
	return
}

func (m *KeySet[K, V]) Values(ctx context.Context, key K) (err error) {
	_, err = m.client.Del(ctx, m.SelfKey(key)).Result()
	return
}

func (m *KeySet[K, V]) DelMembers(ctx context.Context, key K, valueIF V) (err error) {
	value, err := m.marshalFunction(valueIF)
	if err != nil {
		return err
	}
	_, err = m.client.SRem(ctx, m.SelfKey(key), value).Result()
	return err
}

func (m *KeySet[K, V]) DelSet(ctx context.Context, key K) (err error) {
	_, err = m.client.Del(ctx, m.SelfKey(key)).Result()
	return
}

func (m *KeySet[K, V]) GetAndUnmarshalMembers(ctx context.Context, key K) (values []V, err error) {
	valueStrList, err := m.client.SMembers(ctx, m.SelfKey(key)).Result()
	if err != nil {
		return nil, err
	}
	values = []V{}
	for _, vStr := range valueStrList {
		value, err1 := generic_conv.Conv[V](vStr, m.unmarshalFunction)
		if err1 != nil {
			xlog.Error(err1)
			return nil, err1
		}
		values = append(values, value)
	}
	return values, nil
}
