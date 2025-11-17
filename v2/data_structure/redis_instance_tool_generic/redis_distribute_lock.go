package redis_instance_tool

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golib/v2/utils"
	"github.com/golib/v2/xContext/loggers/xlog"
	"github.com/redis/go-redis/v9"
	"math"
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
	interval     time.Duration
	maxRetry     int
	pBase        float64
}

func NewLock(key, sep, saltString string, expire int, interval time.Duration, maxRetry int, pBase float64, client redis.Cmdable) (lock *DLock, err error) {
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
		maxRetry,
		pBase,
	}
	if lock.pBase <= 1 {
		lock.pBase = 1
	}
	if lock.MarshalInterface == nil {
		defaultMarshal.Apply(&lock.RedisToolBase)
	}
	return
}

func (m *DLock) SelfKey(key string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(key+m.saltString)))
}

func (m *DLock) randInterval(t time.Duration) time.Duration {
	step := t + (t/100)*time.Duration(rand.Int()%10+10)
	return step
}

func (m *DLock) RetryScale(i int) (scale float64) {
	scale = math.Pow(m.pBase, float64(i))
	scale = math.Min(scale, 10)
	scale = math.Max(scale, 1)
	return scale
}

func (m *DLock) Lock(ctx context.Context, key string) (ulk *UnlockHdl, err error) {
	lKey := key
	dKey := ""
	for retry := 0; retry < m.maxRetry; retry++ {
		step := (m.randInterval(m.interval) / 1000) * time.Duration(m.RetryScale(retry)*1000)
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
