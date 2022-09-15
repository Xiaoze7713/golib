package redis_instance_tool

import (
	"context"
	"git.singularity-ai.com/backend/library/utils"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

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

func TestDQ(t *testing.T) {
	xlog.SetupLogDefault()
	cli := redis.NewClient(&redis.Options{
		Addr:     "39.99.156.201:6379",
		Password: "",
		DB:       0,
	})
	dq, err := NewDelayQueue(cli)
	type Msg struct {
		Text string `json:"text"`
	}
	ctx := context.Background()
	user := "user"
	err = dq.Add(ctx, user, &Msg{Text: "hi"}, 5, 0)
	if err != nil {
		xlog.Errorf("%v", err)
	}
	time.Sleep(time.Second * 7)
	resList, err := dq.PopL(ctx, 0, utils.UTimeMs())
	if err != nil {
		xlog.Errorf("%v", err)
	}
	xlog.Infof("data %v", resList)
	time.Sleep(time.Second * 5)
}
