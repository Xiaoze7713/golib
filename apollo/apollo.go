/**
 * @Author: wenliangzhang
 * @Description:
 * @File: apollo
 * @Version: 1.0.0
 * @Date: 2022/5/13 6:20 PM
 */
package apollo

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/zouyx/agollo/v4"
	"github.com/zouyx/agollo/v4/env/config"

	"git.singularity-ai.com/backend/library/log"
)

var client *Client

type Client struct {
	*agollo.Client
}

type Config struct {
	Name      string `toml:"name"`
	AppID     string `toml:"app_id"`
	Cluster   string `toml:"cluster"`
	MetaIP    string `toml:"meta_ip"`
	NameSpace string `toml:"namespace"`
}

var defaultTokenConfigPath = "conf/service/apollo.toml"

func Init(filePath string) {
	if filePath == "" {
		filePath = defaultTokenConfigPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Infoln("apollo.toml not exist")
		return
	}
	var c Config
	if _, err := toml.DecodeFile(filePath, &c); err != nil {
		log.Errorf("Can't load config file, %s", err.Error())
		return
	}
	client, _ = NewApolloClient(c)
}

func NewApolloClient(conf Config) (*Client, error) {
	appConf := &config.AppConfig{
		AppID:          conf.AppID,
		Cluster:        conf.Cluster,
		IP:             conf.MetaIP,
		NamespaceName:  conf.NameSpace,
		IsBackupConfig: true,
	}

	apolloClient, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return appConf, nil
	})
	if err != nil {
		log.Errorf("start apollo failed,err=%s", err.Error())
		return nil, err
	}
	return &Client{apolloClient}, nil
}

func GetApolloClient() *Client {
	if client == nil {
		log.Errorln("ApolloClient is nil")
	}
	return client
}
