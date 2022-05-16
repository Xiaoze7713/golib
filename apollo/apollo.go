/**
 * @Author: wenliangzhang
 * @Description:
 * @File: apollo
 * @Version: 1.0.0
 * @Date: 2022/5/13 6:20 PM
 */
package apollo

import (
	"git.singularity-ai.com/backend/library/log"
	"github.com/zouyx/agollo/v4"
	"github.com/zouyx/agollo/v4/env/config"
	"strings"
)

var client *Client

type Client struct {
	*agollo.Client
}

func NewApolloClient(namespaces []string) (*Client, error) {
	conf := &config.AppConfig{
		AppID:          "aiyou",
		Cluster:        "default",
		IP:             "https://config-center-apollo.singularity-ai.com",
		NamespaceName:  strings.Join(namespaces, ","),
		IsBackupConfig: true,
	}

	apolloClient, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return conf, nil
	})
	if err != nil {
		log.Errorf("start apollo failed,err=%s", err.Error())
	}
	client = &Client{apolloClient}
	return client, err
}

func GetApolloClient() *Client {
	if client == nil {
		log.Errorln("ApolloClient is nil")
	}
	return client
}
