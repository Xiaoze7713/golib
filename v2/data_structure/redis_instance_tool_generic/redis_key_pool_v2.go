package redis_instance_tool

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golib/v2/generic_conv"
	"github.com/golib/v2/type_def"
	"github.com/golib/v2/xContext/loggers/xlog"
	"github.com/redis/go-redis/v9"
)

type KeyPoolV2[K type_def.BaseValueType, V any] struct {
	RedisToolBase
	nx     bool
	expire int // second
}

func NewKeyPoolV2[K type_def.BaseValueType, V any](businessKey, sep string, client redis.Cmdable, nx bool, duration int) (keyPool *KeyPoolV2[K, V], err error) {
	if businessKey == "" {
		xlog.Warn("key pool null business key")
	}
	if client == nil {
		err = errors.New("redis pool is nil")
	}
	keyPool = &KeyPoolV2[K, V]{
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

func (m *KeyPoolV2[K, V]) SelfKey(s K) (key string) {
	return m.Prefix() + fmt.Sprintf("%v", s)
}

func (m *KeyPoolV2[K, V]) SelfSep() (sep string) {
	return m.sep
}

func (m *KeyPoolV2[K, V]) ZKey() (key string) {
	return m.Prefix() + "_sort_list"
}

func (m *KeyPoolV2[K, V]) Set(ctx context.Context, key K, valueIF V) (err error) {
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
	_, err = m.client.ZAdd(ctx, m.ZKey(), []redis.Z{{
		Member: m.SelfKey(key),
		Score:  float64(time.Now().UnixMilli()),
	},
	}...).Result()
	if err != nil {
		return err
	}
	return nil
}

func (m *KeyPoolV2[K, V]) Flush(ctx context.Context) (err error) {
	_, err = m.client.ZRemRangeByScore(ctx, m.ZKey(), fmt.Sprintf("0"),
		fmt.Sprintf("%v", time.Now().UnixMilli()-int64(m.expire*1000))).Result()
	if err != nil {
		return
	}
	return nil
}

func (m *KeyPoolV2[K, V]) Count(ctx context.Context) (count int, err error) {
	err = m.Flush(ctx)
	if err != nil {
		return -1, err
	}
	count64, err := m.client.ZCount(ctx, m.ZKey(),
		fmt.Sprintf("0"),
		fmt.Sprintf("%v", time.Now().UnixMilli())).Result()
	if err != nil {
		return -1, err
	}
	return int(count64), nil
}

func (m *KeyPoolV2[K, V]) Get(ctx context.Context, key K) (resStr string, err error) {
	resStr, err = m.client.Get(ctx, m.SelfKey(key)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return
}

func (m *KeyPoolV2[K, V]) Del(ctx context.Context, key K) (err error) {
	_, err = m.client.Del(ctx, m.SelfKey(key)).Result()
	if err != nil {
		return err
	}
	_, err = m.client.ZRem(ctx, m.ZKey(), m.SelfKey(key)).Result()
	return
}

func (m *KeyPoolV2[K, V]) GetAndUnmarshal(ctx context.Context, key K) (value V, err error) {
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
