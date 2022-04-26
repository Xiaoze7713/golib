/**
 * @Author: wenliangzhang
 * @Description:
 * @File: init
 * @Version: 1.0.0
 * @Date: 2022/4/26 10:55 AM
 */
package ali_oss

import (
	"git.singularity-ai.com/backend/ws_service/library/log"
	"github.com/BurntSushi/toml"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

type Config struct {
	AccessKeyID      string `toml:"access_key_id"`
	AccessKeySecret  string `toml:"access_key_secret"`
	EndPointInternal string `toml:"end_point_internal"`
	EndPoint         string `toml:"end_point"`
	DomainName       string `toml:"domain_name"`
	BucketName       string `toml:"bucket_name"`
}

var c *Config

// Init 根据配置文件初始化
func Init(configFile string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Infoln("ali_oss conf file not exist")
		return err
	}
	c = &Config{}
	if _, err := toml.DecodeFile(configFile, c); err != nil {
		log.Errorln("decode conf file not exist")
		return err
	}
	return nil
}

func Client(configFile string) (*oss.Client, error) {
	if c == nil {
		if err := Init(configFile); err != nil {
			return nil, err
		}
	}
	client, err := oss.New(c.EndPoint, c.AccessKeyID, c.AccessKeySecret)
	if err != nil {
		log.Errorf("new oss client failed,err=%s", err.Error())
		return nil, err
	}
	return client, nil
}
