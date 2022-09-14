package apollo_v2

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"sync"
)

var conf *config.AppConfig
var once = sync.Once{}

type ApolloHandler struct {
	nsMap *sync.Map
	conf  *config.AppConfig
}

var Handler *ApolloHandler

func (m *ApolloHandler) GetJsonData(ns string, i interface{}) (err error) {
	_, ok := m.nsMap.Load(ns)
	if !ok {
		err = m.Register([]string{ns})
	}
	cliIfc, _ := m.nsMap.Load(ns)
	cli, ok := cliIfc.(agollo.Client)
	if jsonCfg := cli.GetConfig(ns); jsonCfg != nil {
		content := jsonCfg.GetValue("content")
		if content != "" {
			fmt.Println(content)
			err = jsoniter.UnmarshalFromString(content, i)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *ApolloHandler) Register(nsList []string) (err error) {
	newNsList := []string{}
	for _, ns := range nsList {
		_, exist := m.nsMap.Load(ns)
		if exist {
			continue
		} else {
			newNsList = append(newNsList, ns)
		}
	}
	if len(newNsList) > 0 {
		cli, err := m.newApolloClient(newNsList)
		if err != nil {
			return err
		}
		for _, ns := range newNsList {
			m.nsMap.Store(ns, cli)
		}
	}
	return nil
}

func Init(host, cluster, appID string, ns []string) (err error) {
	once.Do(func() {
		Handler, err = NewInstance(host, cluster, appID, ns)
	})
	return
}

func NewInstance(host, cluster, appID string, ns []string) (hdl *ApolloHandler, err error) {
	conf = &config.AppConfig{
		AppID:            appID,
		Cluster:          cluster,
		IP:               host,
		IsBackupConfig:   true,
		MustStart:        true,
		BackupConfigPath: "conf/apollo/",
	}
	hdl = &ApolloHandler{
		nsMap: &sync.Map{},
		conf:  nil,
	}
	hdl.conf = conf
	if len(ns) > 0 {
		err = hdl.Register(ns)
	}
	return
}

func (m *ApolloHandler) newApolloClient(ns []string) (agollo.Client, error) {
	apolloClient, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return &config.AppConfig{
			AppID:             m.conf.AppID,
			Cluster:           m.conf.Cluster,
			NamespaceName:     strings.Join(ns, ","),
			IP:                m.conf.IP,
			IsBackupConfig:    m.conf.IsBackupConfig,
			SyncServerTimeout: 30,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return apolloClient, err
}
