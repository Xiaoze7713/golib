package apollo_v2

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

var conf *config.AppConfig

func NewApolloClient(host, cluster, appID string) (agollo.Client, error) {
	conf = &config.AppConfig{
		AppID:   appID,
		Cluster: cluster,
		IP:      host,
		//NamespaceName:  "expression_config",
		IsBackupConfig: true,
	}
	apolloClient, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return conf, nil
	})
	if err != nil {
		return nil, err
		//xlog.Error("start apollo failed,err=%s", err.Error())
	}
	return apolloClient, err
}
