package apollo_v2

import (
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"sync"
)

var conf *config.AppConfig

//var once = sync.Once{}

type ApolloDataType interface {
	any
}

type ApolloCliHandler struct {
	nsMap *sync.Map
	Cli   agollo.Client
	AppID string
}

type SimpleApolloConfig struct {
	Host    string `toml:"host"`
	Cluster string `toml:"cluster"`
}

type AppConfig struct {
	AppID     string `toml:"app_id"`
	Namespace string `toml:"namespace"`
}

func newInstance(host, cluster, appID string, ns []string) (hdl *ApolloCliHandler, err error) {
	conf = &config.AppConfig{
		AppID:             appID,
		Cluster:           cluster,
		IP:                host,
		IsBackupConfig:    true,
		MustStart:         true,
		BackupConfigPath:  "conf/apollo/",
		NamespaceName:     strings.Join(ns, ","),
		SyncServerTimeout: 30,
	}
	apolloClient, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return conf, nil
	})
	if err != nil {
		return nil, err
	}
	hdl = &ApolloCliHandler{
		nsMap: &sync.Map{},
		Cli:   apolloClient,
	}
	return
}

func combineKey(appID, namespace string) string {
	return appID + "##" + namespace
}

type ApolloManager struct {
	appMap  *sync.Map
	dataMap *sync.Map
	conf    *SimpleApolloConfig
}

var Handler *ApolloManager
var once = sync.Once{}

func Init(conf *SimpleApolloConfig) (err error) {
	once.Do(func() {
		Handler = &ApolloManager{
			appMap:  &sync.Map{},
			conf:    conf,
			dataMap: &sync.Map{},
		}
	})
	return
}

//
func (m *ApolloManager) Register(appID string, nsList []string) (err error) {
	newNsList := []string{}
	var apolloHandler *ApolloCliHandler
	app, ok := m.appMap.Load(appID)
	if !ok {
		newNsList = nsList
	} else {
		apolloHandler, ok = app.(*ApolloCliHandler)
		if ok {
			for _, ns := range nsList {
				_, exist := apolloHandler.nsMap.Load(ns)
				if exist {
					continue
				} else {
					newNsList = append(newNsList, ns)
				}
			}
		}
	}
	if len(newNsList) > 0 {
		if apolloHandler != nil {
			for _, ns := range newNsList {
				val := apolloHandler.Cli.GetConfig(ns)
				if val == nil {
					xlog.Errorf("unknown config name app %s, namespace %s in cluster %s, host %s", appID, ns, m.conf.Cluster, m.conf.Host)
				}
				//m.dataMap.Store(combineKey(appID, ns), val)
			}
		} else {
			hdl, err := newInstance(m.conf.Host, m.conf.Cluster, appID, newNsList)
			if err != nil {
				return err
			}
			m.appMap.Store(appID, hdl)
		}
	}
	return nil
}
func (m *ApolloManager) getJsonData(appID, ns string, i interface{}) (err error) {
	failedKey := fmt.Sprintf("config name app %s, namespace %s", appID, ns)
	_, ok := m.appMap.Load(appID)
	if !ok {
		err = m.Register(appID, []string{ns})
		if err != nil {
			return err
		}
	}
	cliIfc, _ := m.appMap.Load(appID)
	handler, ok := cliIfc.(*ApolloCliHandler)
	if !ok {
		return errors.New("failed get cli " + failedKey)
	}
	if jsonCfg := handler.Cli.GetConfig(ns); jsonCfg != nil {
		content := jsonCfg.GetValue("content")
		if content != "" {
			//fmt.Println(content)
			err = jsoniter.UnmarshalFromString(content, i)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("failed get config " + failedKey)
}

func (m *ApolloManager) cacheData(appID, ns string, data interface{}) {
	m.dataMap.Store(combineKey(appID, ns), data)
}

func GetData[T ApolloDataType](appID, ns string, t *T) (err error) {
	err = Handler.getJsonData(appID, ns, t)
	if err == nil {
		Handler.cacheData(appID, ns, t)
		return
	}
	data, ok := Handler.dataMap.Load(combineKey(appID, ns))
	if ok {
		t, ok = data.(*T)
		if ok {
			xlog.Warnf("get config failed , use cache")
			return nil
		}
	}
	xlog.Error("get config error %v", err)
	return err
}
