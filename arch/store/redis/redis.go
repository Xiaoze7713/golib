/**
 * @Author: wenliangzhang
 * @Description:
 * @File: redis
 * @Version: 1.0.0
 * @Date: 2022/4/18 2:29 PM
 */
package redis

import (
	rds "github.com/gomodule/redigo/redis"

	"git.singularity-ai.com/backend/library/log"
)

func Get(key string) (string, error) {
	return sampleDoString("GET", key)
}

func Set(key, value string, ttl int64) error {
	var err error
	if ttl > 0 {
		_, err = sampleDo("SET", key, value, "EX", ttl)
	} else {
		_, err = sampleDo("SET", key, value)
	}
	return err
}

func SetNx(key, value string, ttl int64) (interface{}, error) {
	var err error
	var resp interface{}
	if ttl > 0 {
		resp, err = sampleDo("SET", key, value, "EX", ttl, "NX")
	} else {
		resp, err = sampleDo("SET", key, value, "NX")
	}
	return resp, err
}

func Del(key string) error {
	_, err := sampleDo("DEL", key)
	if err != nil {
		log.Fatalf("Del err=%v", err)
	}
	return err
}

func Incr(key string) (int64, error) {
	return sampleDoInt64("INCR", key)
}

func Incrby(key string, num int) (int64, error) {
	return sampleDoInt64("INCRBY", key, num)
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

	value, err := sampleDo(method, params...)
	if err != nil {
		return "", err
	}

	return value, nil
}

func Rpush(queue, value string) error {
	_, err := sampleDo("RPUSH", queue, value)
	return err
}

func Lpop(queue string) (string, error) {
	return sampleDoString("LPOP", queue)
}

func HSet(key, field, value string) error {
	var err error
	_, err = sampleDo("HSET", key, field, value)
	return err
}

func HGet(key, field string) (string, error) {
	return sampleDoString("HGET", key, field)
}

func HDel(key string, field string) (string, error) {
	return sampleDoString("HDEL", key, field)
}

func SAdd(key string, data []string) error {
	var params []interface{}
	params = append(params, key)
	for _, value := range data {
		params = append(params, value)
	}
	_, err := sampleDo("SADD", params...)
	return err
}

func SIsMember(key string, value string) (int, error) {
	return sampleDoInt("SISMEMBER", key, value)
}

func SMembers(key string) ([]string, error) {
	return rds.Strings(sampleDo("SMEMBERS", key))
}

func SRem(key string, data []string) error {
	var params []interface{}
	params = append(params, key)
	for _, value := range data {
		params = append(params, value)
	}
	_, err := sampleDo("SREM", params...)
	return err
}
