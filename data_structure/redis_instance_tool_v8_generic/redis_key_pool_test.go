package redis_instance_tool

import (
	"context"
	"fmt"
	"git.singularity-ai.com/backend/library/v2/xContext"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func TestTokenPool(t *testing.T) {
	//cmdAble := redis.NewClient(&redis.Options{
	//	Addr:         redisCfg.Host,
	//	Password:     redisCfg.Password,
	//	Username:     redisCfg.UserName,
	//	DB:           int(redisCfg.DBNumber),
	//	DialTimeout:  time.Second * 3,
	//	ReadTimeout:  time.Second * 3,
	//	WriteTimeout: time.Second * 3,
	//})
	//rdsPool.Close()
	//kpl := TokenPool{KeyPool{RedisCli: RedisCli}}
	//err := kpl.Set("123", "123", 30, true)
	//fmt.Printf("%v\n", err)
	//res, err := kpl.Get("123")
	//fmt.Printf("%v %v\n", res, err)
}

func TestCtx(t *testing.T) {
	bc_ctx := context.Background()
	fmt.Printf("%+v\n", &bc_ctx)
	valueCtx := context.WithValue(context.Background(), "hi", "mingze")
	fmt.Printf("%+v\n", &valueCtx)
	fmt.Printf("%+v\n", valueCtx.Value("hi"))
	newVCtx, cancel := context.WithCancel(valueCtx)
	fmt.Printf("%+v\n", newVCtx.Value("hi"))
	fmt.Printf("%+v\n", context.Background().Value("hi"))

	defer cancel()
}

func TestJobKey(t *testing.T) {

}

type XXX struct {
	TimeMs int64  `json:"time_ms"`
	Value  string `json:"value"`
}

func TestExpirePool(t *testing.T) {
	xContext.InitByNullOpt()
	defer xlog.Close()
	defer time.Sleep(time.Second * 3)
	ctx := xContext.NewXContext("xx")
	defer ctx.Fin()
	cmdAble := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		DB:           0,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	})
	defer cmdAble.Close()
	zkp, err := NewKeyPoolV2[string, *XXX]("keypool_test", "_", cmdAble, false, 30)
	if err != nil {
		ctx.Error(err)
		return
	}
	go func() {
		for i := 0; i < 1800; i++ {
			key := fmt.Sprintf("%v", i)
			xx := &XXX{
				TimeMs: time.Now().UnixMilli(),
				Value:  key,
			}
			err = zkp.Set(ctx, key, xx)
			if err != nil {
				xlog.Error(err)
			}
			if i < 120 {
				time.Sleep(time.Millisecond * 1000)
			} else {
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()
	go func() {
		for i := 0; i < 1800; i++ {
			//key := fmt.Sprintf("%v", i-1)
			//xx := &XXX{
			//	TimeMs: time.Now().UnixMilli(),
			//	Value:  key,
			//}
			//err = zkp.Del(ctx, key)
			//if err != nil {
			//	xlog.Error(err)
			//}
			if i < 120 {
				time.Sleep(time.Millisecond * 1000)
			} else {
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()
	for i := 0; i < 1800; i++ {
		time.Sleep(time.Second * 1)
		cnt, err1 := zkp.Count(ctx)
		if err1 != nil {
			xlog.Error(err1)
		} else {
			xlog.Info(cnt)
		}
	}
}
