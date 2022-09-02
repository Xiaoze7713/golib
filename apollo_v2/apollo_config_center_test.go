package apollo_v2

import (
	"fmt"
	"git.singularity-ai.com/backend/library/utils"
	"testing"
)

func Test_apollo_run(t *testing.T) {
	simpleConfig := &SimpleConfig{
		Cluster: "default",
		Host:    "https://config-center-apollo.singularity-ai.com",
		AppID:   "expression_v1",
	}
	xx, err := NewApolloClient(simpleConfig.Host, simpleConfig.Cluster, simpleConfig.AppID)
	if err != nil {
		fmt.Println(err)
		return
	}
	ns := "expression_config.json"
	//mm := map[string]interface{}{}
	_ = xx.GetConfigAndInit(ns)
	//cache := xx.GetApolloConfigCache()
	res := xx.GetConfig(ns)
	fmt.Println(utils.MustJson(res))
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
