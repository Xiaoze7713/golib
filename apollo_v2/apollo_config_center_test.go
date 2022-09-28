package apollo_v2

import (
	"git.singularity-ai.com/backend/library/utils"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"testing"
	"time"
)

func Test_apollo_run(t *testing.T) {
	xlog.SetupLogDefault()
	defer func() {
		time.Sleep(time.Second * 5)
	}()
	simpleConfig := &SimpleApolloConfig{
		Cluster: "default",
		Host:    "https://apollo-dev.singularity-ai.com",
	}
	AppID := "expression_v1"
	nsList := []string{"expression_config.json"}
	//nsList := []string{}
	err := Init(simpleConfig)
	if err != nil {
		xlog.Error(err)
		return
	}
	err = Handler.Register(AppID, nsList)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Info("init finish")
	ns := "expression_config.json"
	//mm := map[string]interface{}{}
	//_ = xx.GetConfig(ns)
	//cache := xx.GetApolloConfigCache()
	mm := map[string]interface{}{}
	err = GetData[map[string]interface{}](AppID, ns, &mm)
	if err != nil {
		xlog.Error(err)
	}
	xlog.Info(utils.MustJson(mm))
	//data, _ := conf.Get("content")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(utils.MustString(data))
	// fmt.Println(data)
	// err = jsoniter.UnmarshalFromString(data.(string), &mm)
	//err = jsoniter.UnmarshalFromString(data, &mm)
	//fmt.Println(err)
}
