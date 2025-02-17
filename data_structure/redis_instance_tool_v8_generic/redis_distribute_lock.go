package redis_instance_tool

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golib/v2/utils"
	"github.com/golib/v2/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

var LockBusy = errors.New("lock busy")

type UnlockHdl struct {
	dLock *DLock
	dKey  string
	key   string
	ctx   context.Context
}

func (m *UnlockHdl) Unlock() (success bool, err error) {
	success, err = m.dLock.unLock(m.ctx, m.key, m.dKey)
	if err != nil {
		return false, err
	}
	return true, nil
}

type DLock struct {
	RedisToolBase
	//RedisCli *redis.Client
	//selfKey  string
	//sep      string
	expireSecond int // second
	saltString   string
	maxInterval  time.Duration
}

func NewLock(key, sep, saltString string, expire int, interval time.Duration, client redis.Cmdable) (lock *DLock, err error) {
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
		interval,
	}
	if lock.MarshalInterface == nil {
		defaultMarshal.Apply(&lock.RedisToolBase)
	}
	return
}

func (m *DLock) SelfKey(key string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(key+m.saltString)))
}

func (m *DLock) NewLock(ctx context.Context, key string) (ulk *UnlockHdl, err error) {
	lKey := key
	step := time.Millisecond * 25
	dKey := ""
	for i := time.Duration(0); i < m.maxInterval; i += step {
		dKey, err = m.lock(ctx, key)
		if err != nil {
			return nil, err
		}
		if dKey != "" {
			break
		}
		time.Sleep(step)
	}
	if dKey == "" {
		return nil, LockBusy
	}
	ulk = &UnlockHdl{
		dLock: m,
		dKey:  dKey,
		key:   lKey,
		ctx:   ctx,
	}
	return ulk, nil
}

func (m *DLock) lock(ctx context.Context, key string) (dlKey string, err error) {
	dlKey, err = utils.StringsMd5([]string{time.Now().String(), fmt.Sprintf("%f", rand.Float64())})
	if err != nil {
		return "", err
	}
	res, err := m.client.SetNX(ctx, m.SelfKey(key), dlKey, time.Duration(m.expireSecond)*time.Second).Result()
	if err != nil {
		return "", err
	}
	if res {
		return dlKey, nil
	} else {
		return "", nil
	}
}

func (m *DLock) unLock(ctx context.Context, key string, dlKey string) (success bool, err error) {
	//v := fmt.Sprintf("%d", time.Now().UnixMilli())
	s := `local val = redis.call("GET",KEYS[1]);
if val == ARGV[1]
then
	return redis.call('DEL', KEYS[1])
else
	return 0
end`
	res, err := m.client.Eval(ctx, s, []string{m.SelfKey(key)}, dlKey).Result()
	if err != nil {
		return false, err
	}
	resStr, ok := res.(int64)
	if !ok {
		return false, errors.New("eval lua script result err")
	}
	if resStr == 1 {
		return true, nil
	}
	return false, nil
}
