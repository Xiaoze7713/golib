package redis_instance_tool

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/xContext"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"time"
)

type DLock struct {
	RedisToolBase
	//RedisCli *redis.Client
	//selfKey  string
	//sep      string
	expire     int // second
	saltString string
}

func NewLock(key, sep, saltString string, expire int, client redis.Cmdable) (lock *DLock, err error) {
	if key == "" {
		xlog.Warn("counter null business key")
		//err = errors.New("null business key")
	}
	if client == nil {
		err = errors.New("redis pool is nil")
	}
	lock = &DLock{
		RedisToolBase{
			client:  client,
			selfKey: key,
			sep:     sep,
		},
		expire,
		saltString,
	}
	if lock.MarshalInterface == nil {
		defaultMarshal.Apply(&lock.RedisToolBase)
	}
	return
}

func (m *DLock) SelfKey(key string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(key+m.saltString)))
}

func (m *DLock) Lock(ctx *xContext.XContext, key string) (err error) {
	v := fmt.Sprintf("%d", time.Now().UnixMilli())
	_, err = m.client.Set(ctx, m.SelfKey(key), v, time.Duration(m.expire)*time.Second).Result()
	if err != nil {
		return err
	}
	return nil
}

func (m *DLock) UnLock(ctx *xContext.XContext, key string) (err error) {
	//v := fmt.Sprintf("%d", time.Now().UnixMilli())
	_, err = m.client.Del(ctx, m.SelfKey(key)).Result()
	if err != nil {
		return err
	}
	return nil
}
