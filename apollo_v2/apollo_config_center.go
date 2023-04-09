package apollo_v2

import (
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/v2/utils"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
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

func newInstance(host, cluster, appID string, nsList []string) (hdl *ApolloCliHandler, err error) {
	ns := strings.Join(nsList, ",")
	apolloClient, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		conf = &config.AppConfig{
			AppID:          appID,
			Cluster:        cluster,
			IP:             host,
			IsBackupConfig: false,
			MustStart:      true,
			//BackupConfigPath:  "./conf/apollo/",
			NamespaceName:     ns,
			SyncServerTimeout: 30,
		}
		return conf, nil
	})
	if err != nil {
		return nil, err
	}
	hdl = &ApolloCliHandler{
		nsMap: &sync.Map{},
		Cli:   apolloClient,
		AppID: appID,
	}
	for _, ns = range nsList {
		hdl.nsMap.Store(ns, true)
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
	lock    *sync.RWMutex
}

var Handler *ApolloManager
var once = sync.Once{}

func Init(conf *SimpleApolloConfig) (err error) {
	once.Do(func() {
		Handler = &ApolloManager{
			appMap:  &sync.Map{},
			conf:    conf,
			dataMap: &sync.Map{},
			lock:    &sync.RWMutex{},
		}
	})
	return
}

//
func (m *ApolloManager) Register(appID string, nsList []string) (err error) {
	newNsList := []string{}
	var apolloHandler *ApolloCliHandler
	m.lock.Lock()
	xlog.Infof("new app %v namespace %v", appID, utils.MustJson(nsList))
	defer m.lock.Unlock()
	app, ok := m.appMap.Load(appID)
	if !ok {
		newNsList = nsList
	} else {
		apolloHandler, ok = app.(*ApolloCliHandler)
		if !ok {
			xlog.Errorf("handler not found app %s %v", appID, utils.MustJson(nsList))
		}
	}
	if apolloHandler == nil {
		apolloHandler, err = newInstance(m.conf.Host, m.conf.Cluster, appID, newNsList)
		if err != nil {
			xlog.Error(err)
			return err
		}

	} else {
		newNsList = nsList
		apolloHandler.nsMap.Range(func(key, value any) bool {
			ns := key.(string)
			println(ns)
			newNsList = append(newNsList, ns)
			return true
		})
		println("new ns ", utils.MustJson(newNsList))
		apolloHandler, err = newInstance(m.conf.Host, m.conf.Cluster, appID, newNsList)
	}
	if len(newNsList) > 0 {
		for _, ns := range newNsList {
			val := apolloHandler.Cli.GetConfig(ns)
			if val == nil {
				xlog.Errorf("unknown config name app %s, namespace %s in cluster %s, host %s", appID, ns, m.conf.Cluster, m.conf.Host)
			}
			//m.dataMap.Store(combineKey(appID, ns), val)
		}
	}
	m.appMap.Store(appID, apolloHandler)
	return nil
}
func (m *ApolloManager) getCliHandler(appID string) (*ApolloCliHandler, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	cliIfc, ok := m.appMap.Load(appID)
	if !ok {
		return nil, false
	} else {
		handler, ok := cliIfc.(*ApolloCliHandler)
		if !ok {
			return nil, false
		}
		return handler, true
	}
}

func (m *ApolloManager) getJsonData(appID, ns string, i interface{}) (err error) {
	failedKey := fmt.Sprintf("config name app %s, namespace %s", appID, ns)
	handler, ok := m.getCliHandler(appID)
	if !ok {
		err = m.Register(appID, []string{ns})
		if err != nil {
			return err
		}
		handler, ok = m.getCliHandler(appID)
		if !ok {
			return errors.New("get handler failed " + failedKey)
		}
	}
	jsonCfg := handler.Cli.GetConfigAndInit(ns)
	if jsonCfg == nil {
		err = m.Register(appID, []string{ns})
		if err != nil {
			return err
		}
		handler, ok = m.getCliHandler(appID)
		if !ok {
			return errors.New("get handler failed " + failedKey)
		}
	}
	jsonCfg = handler.Cli.GetConfigAndInit(ns)
	if jsonCfg != nil {
		content := jsonCfg.GetValue("content")
		if content != "" {
			//fmt.Println(content)
			err = jsoniter.UnmarshalFromString(content, i)
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("null content " + failedKey)
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
	xlog.Error(err)
	return err
}
