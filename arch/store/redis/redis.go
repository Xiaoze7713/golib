/**
 * @Author: wenliangzhang
 * @Description:
 * @File: redis
 * @Version: 1.0.0
 * @Date: 2022/4/18 2:29 PM
 */
package redis

import (
	"fmt"
	"git.singularity-ai.com/backend/library/log"
	"github.com/BurntSushi/toml"
	rds "github.com/gomodule/redigo/redis"
	"os"
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

func Get(key string) (string, error) {

	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close conn err=%v", err)
		}
	}(conn)
	value, err := conn.Do("GET", key)
	if err != nil || value == nil {
		return "", err
	}
	return string(value.([]uint8)), nil
}

func Set(key, value string, ttl int64) error {

	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("close conn err=%v", err)
		}
	}(conn)
	var err error
	if ttl > 0 {
		_, err = conn.Do("SET", key, value, "EX", ttl)
	} else {
		_, err = conn.Do("SET", key, value)
	}
	return err
}

func Del(key string) error {
	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close conn err=%v", err)
		}
	}(conn)
	_, err := conn.Do("DEL", key)
	if err != nil {
		log.Fatalf("Del err=%v", err)
	}
	return err
}

func Incr(key string) (int64, error) {
	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close conn err=%v", err)
		}
	}(conn)
	reply, err := conn.Do("INCR", key)
	if err != nil {
		log.Fatalf("Incr err=%v", err)
	}
	return reply.(int64), err
}

func Incrby(key string) (int64, error) {
	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close conn err=%v", err)
		}
	}(conn)
	reply, err := conn.Do("INCRBY", key)
	if err != nil {
		log.Fatalf("Incrby err=%v", err)
	}
	return reply.(int64), err
}

// Expire
// 当key的值为数字时 返回value为int64,
// 当key的值非数字时，返回value为string
func Expire(key string, ttl int) (interface{}, error) {
	var params []interface{}
	var method string

	if ttl > 0 {
		method = "EXPIRE"
		params = append(params, key, ttl)
	} else {
		method = "PERSIST"
		params = append(params, ttl)
	}

	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close conn err=%v", err)
		}
	}(conn)

	value, err := conn.Do(method, params...)
	if err != nil {
		return "", err
	}

	return value, nil
}

func Rpush(queue, value string) error {
	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close conn err=%v", err)
		}
	}(conn)
	_, err := conn.Do("RPUSH", queue, value)
	return err
}

func Lpop(queue string) (string, error) {
	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close conn err=%v", err)
		}
	}(conn)
	value, err := conn.Do("LPOP", queue)
	if err != nil || value == nil {
		return "", err
	}
	return string(value.([]uint8)), nil
}
