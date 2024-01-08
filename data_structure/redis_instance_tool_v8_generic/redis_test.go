package redis_instance_tool

import (
	"context"
	"fmt"
	"github.com/Xiaoze7713/golib/v3/utils"
	"github.com/Xiaoze7713/golib/v3/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"reflect"
	"sync"
	"testing"
	"time"
)

func easyCli(env string) (c redis.Cmdable) {
	if env == "dev" || env == "test" {
		c = redis.NewClient(&redis.Options{
			Addr:     "39.99.156.201:6379",
			Password: "",
			DB:       0,
		})
		return
	} else if env == "local" {
		c = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})
		return
	}
	return nil
}

func TestTable(t *testing.T) {
	xlog.SetupLogDefault()
	cli := redis.NewClient(&redis.Options{
		Addr:     "39.99.156.201:6379",
		Password: "",
		DB:       0,
	})
	tb, err := NewTable("", " ", cli)
	tableKey := "001"
	ctx := context.Background()
	err = tb.SetColumn(ctx, tableKey, "userName", "小泽")
	if err != nil {
		xlog.Errorf("%v", err)
	}
	err = tb.SetColumns(ctx, tableKey, map[string]interface{}{"userName": "小白", "robotName": "小奇"})
	if err != nil {
		xlog.Errorf("%v", err)
	}
	valueMap, err := tb.GetColumns(ctx, tableKey, "userName", "robotName")
	if err != nil {
		xlog.Errorf("%v", err)
	}
	xlog.Debugf("%v", valueMap)
	for k, v := range valueMap {
		xlog.Debugf("%v : %v", k, utils.UnMarshalToString(v))
	}
	time.Sleep(time.Second * 3)
	return
}
func TestCounter(t *testing.T) {
	xlog.SetupLogDefault()
	cli := redis.NewClient(&redis.Options{
		Addr:     "39.99.156.201:6379",
		Password: "",
		DB:       0,
	})
	ctr, err := NewCounter("", "", cli)
	tableKey := "001"
	ctx := context.Background()
	v, err := ctr.Add(ctx, tableKey, 10)
	if err != nil {
		xlog.Errorf("%v", err)
	}
	xlog.Infof("%v", v)
	v, err = ctr.Reduce(ctx, tableKey, 3)
	if err != nil {
		xlog.Errorf("%v", err)
	}
	xlog.Infof("%v", v)
	v, err = ctr.Get(ctx, tableKey)
	if err != nil {
		xlog.Errorf("%v", err)
	}
	xlog.Infof("%v", v)
	time.Sleep(time.Second * 3)
	return
}

type XX struct {
	A string `json:"a"`
	B string `json:"b"`
}

func TestDQ(t *testing.T) {
	xlog.SetupLogDefault()
	cli := redis.NewClient(&redis.Options{
		Addr:     "39.99.233.6:6379",
		Password: "redis",
		DB:       0,
	})
	dq, err := NewDelayQueue[*XX]("dq_test", "_", cli)
	type Msg struct {
		Text string `json:"text"`
	}
	ctx := context.Background()
	user := "user"
	err = dq.Add(ctx, user, &XX{
		A: "xa",
		B: "xb",
	}, 5, 10)
	if err != nil {
		xlog.Errorf("%v", err)
	}
	time.Sleep(time.Second * 7)
	resList, err := dq.PopL(ctx, 0, utils.UTimeMs())
	if err != nil {
		xlog.Errorf("%v", err)
	}
	xlog.Infof("data %v", utils.MustJson(resList))

	time.Sleep(time.Second * 5)
}

func TestNewKeyPool(t *testing.T) {
	xlog.SetupLogDefault()
	cli := redis.NewClient(&redis.Options{
		Addr:     "39.99.233.6:6379",
		Password: "redis",
		DB:       0,
	})
	pool, err := NewKeyPool[string, *XX]("dq_test", "_", cli, false, 10)
	ctx := context.Background()
	user := "user"
	err = pool.Set(ctx, user, &XX{
		A: "xa",
		B: "xb",
	})
	if err != nil {
		xlog.Errorf("%v", err)
	}
	time.Sleep(time.Second * 7)
	res, err := pool.GetAndUnmarshal(ctx, user)
	if err != nil {
		xlog.Errorf("%v", err)
	}
	xlog.Infof("data %v", utils.MustJson(res))

	time.Sleep(time.Second * 5)
}

func TestGeneric(t *testing.T) {
	var tt int = 10
	v := reflect.New(reflect.TypeOf(tt))
	println(*(v.Interface().(*int)))
	*(v.Interface().(*int)) = 10
	println(*(v.Interface().(*int)))
}

func TestDLock_Lock(t *testing.T) {
	xlog.SetupLogDefault()
	cli := redis.NewClient(&redis.Options{
		Addr:     "39.99.233.6:6379",
		Password: "redis",
		DB:       0,
	})
	for i := 0; i < 100; i++ {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool, err := NewLock("test", "-", "sai", 3, time.Millisecond*1500, cli)
			ctx := context.Background()
			key01 := "key_01"
			ulk, err := pool.NewLock(ctx, key01)
			if err != nil {
				xlog.Errorf("%v", err)
				return
			} else if ulk != nil {
				xlog.Infof("lock success")
			} else {
				xlog.Infof("lock failed")
				return
			}
			time.Sleep(time.Millisecond * 10)
			success, err := pool.unLock(ctx, key01, "123")
			if err != nil {
				xlog.Error(err)
				return
			}
			xlog.Infof("try unlock %v", success)
			success, err = ulk.Unlock()
			if err != nil {
				xlog.Error(err)
				return
			}
			xlog.Infof("try unlock 2 %v", success)
		}()
		wg.Wait()
	}

	time.Sleep(time.Second * 5)
}

func TestNewQueue(t *testing.T) {
	xlog.SetupLogDefault()
	cli := redis.NewClient(&redis.Options{
		Addr:     "39.99.233.6:6379",
		Password: "redis",
		DB:       0,
	})
	q, err := NewQueue[string]("test", cli)
	if err != nil {
		xlog.Error(err)
		t.Fail()
	}
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			val := fmt.Sprintf("%06d", i)
			err = q.Push(ctx, val)
			if err != nil {
				xlog.Error(err)
				continue
			}
			qLen, err := q.Len(ctx)
			if err != nil {
				xlog.Error(err)
				continue
			}
			xlog.Info(qLen)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for {
			s, err := q.Pop(ctx, true)
			if err != nil {
				xlog.Error(err)
			}
			xlog.Info(s)
		}
	}()
	wg.Wait()
}

type TestData struct {
	ID      int    `json:"id"`
	RandStr string `json:"rand_str"`
}

func TestKeySet(t *testing.T) {
	cli := easyCli("local")
	xlog.SetupLogDefault()
	defer xlog.Close()
	defer time.Sleep(time.Second * 3)
	keySet, _ := NewKeySet[string, *TestData]("test", "-", cli)
	uuid := "user10"
	ctx := context.Background()
	cnt, err := keySet.Count(ctx, uuid)
	if err != nil {
		xlog.Error(err)
		return
	}
	dataList := []*TestData{}
	xlog.Infof("cnt %v", cnt)
	for i := 0; i < 5; i++ {
		newData := &TestData{
			ID:      i,
			RandStr: fmt.Sprintf("%v", i),
		}
		err = keySet.Add(ctx, uuid, newData)
		if err != nil {
			xlog.Error(err)
			return
		}
		dataList = append(dataList, newData)
	}
	cnt, err = keySet.Count(ctx, uuid)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Infof("cnt %v", cnt)
	err = keySet.DelMembers(ctx, uuid, dataList[3])
	if err != nil {
		xlog.Error(err)
		return
	}
	cnt, err = keySet.Count(ctx, uuid)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Infof("cnt %v", cnt)
	ok, err := keySet.ExistMember(ctx, uuid, dataList[1])
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Infof("exist %v", ok)
	ok, err = keySet.ExistMember(ctx, uuid, dataList[3])
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Infof("exist %v", ok)
	newList, err := keySet.GetAndUnmarshalMembers(ctx, uuid)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Infof("members %v", utils.MustJson(newList))
}
