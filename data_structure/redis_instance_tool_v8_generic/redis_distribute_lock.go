package redis_instance_tool

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/utils"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"math/rand"
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

func (m *DLock) Lock(ctx context.Context, key string) (dlKey string, err error) {
	dlKey, err = utils.StringsMd5([]string{time.Now().String(), fmt.Sprintf("%f", rand.Float64())})
	if err != nil {
		return "", err
	}
	res, err := m.client.SetNX(ctx, m.SelfKey(key), dlKey, time.Duration(m.expire)*time.Second).Result()
	if err != nil {
		return "", err
	}
	if res {
		return dlKey, nil
	} else {
		return "", nil
	}
}

func (m *DLock) UnLock(ctx context.Context, key string, dlKey string) (success bool, err error) {
	//v := fmt.Sprintf("%d", time.Now().UnixMilli())
	s := `local key = KEYS[1]
local val = redis.call("GET", key);
if val == ARGV[1]
then
	redis.call('DEL', KEYS[1])
	return 1
else
	return 0
end`
	res, err := m.client.Eval(ctx, s, []string{key}, dlKey).Result()
	if err != nil {
		return false, err
	}
	if res == "1" {
		return success, nil
	}
	return false, nil
}
