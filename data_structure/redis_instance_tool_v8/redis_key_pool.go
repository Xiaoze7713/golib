package redis_instance_tool

import (
	"context"
	"errors"
	"git.singularity-ai.com/backend/library/utils"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"time"
)

type KeyPool struct {
	RedisToolBase
	nx     bool
	expire int // second
}

func NewKeyPool(businessKey, sep string, client *redis.Client, nx bool, duration int) (keyPool *KeyPool, err error) {
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
	return
}

func (m *KeyPool) SelfKey(s string) (key string) {
	return m.selfKey + m.sep + s
}

func (m *KeyPool) SelfSep() (sep string) {
	return m.sep
}

func (m *KeyPool) Set(ctx context.Context, key string, valueIF interface{}) (err error) {
	//value, ok := valueIF.(string)
	//if !ok {
	value, err := utils.MarshalToString(valueIF)
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

func (m *KeyPool) GetAndUnmarshalJson(ctx context.Context, key string, ifc interface{}, unmarshal func(s string, ifc interface{}) error) (err error) {
	resStr, err := m.Get(ctx, key)
	if err != nil {
		return err
	}
	err = unmarshal(resStr, ifc)
	if err != nil {
		return err
	}
	return nil
}
