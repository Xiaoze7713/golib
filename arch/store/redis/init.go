/**
 * @Author: wenliangzhang
 * @Description:
 * @File: init
 * @Version: 1.0.0
 * @Date: 2022/7/5 8:01 PM
 */
package redis

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	rds "github.com/gomodule/redigo/redis"

	"git.singularity-ai.com/backend/library/v2/log"
)

type Client struct {
	C *rds.Pool
}

type Config struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Password string `toml:"password"`
	DBNumber int64  `toml:"db_no"`
}

var defaultRedisConfigPath = "conf/service/redis.toml"

var client *Client

func Init(filePath string) error {
	var config Config
	if filePath == "" {
		filePath = defaultRedisConfigPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("redis.toml not exist")
		return nil
	}
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		panic(fmt.Sprintf("Can't load config file, %s", err.Error()))
	}
	rdsPool := &rds.Pool{
		MaxIdle:   20,
		MaxActive: 100,
		Dial: func() (rds.Conn, error) {
			c, err := rds.Dial("tcp", config.Host)
			if err != nil {
				return nil, err
			}
			if config.Password != "" {
				if _, err := c.Do("AUTH", config.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if config.DBNumber > 0 {
				if _, err := c.Do("SELECT", config.DBNumber); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}
	client = &Client{
		C: rdsPool,
	}
	return nil
}
