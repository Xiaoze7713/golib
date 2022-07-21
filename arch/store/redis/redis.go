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

	"git.singularity-ai.com/backend/library/arch/web"
)

func Get(ctx *web.WebContext, key string) (string, error) {
	return sampleDoString(ctx, "GET", key)
}

func Set(ctx *web.WebContext, key, value string, ttl int64) error {
	var err error
	if ttl > 0 {
		_, err = sampleDo(ctx, "SET", key, value, "EX", ttl)
	} else {
		_, err = sampleDo(ctx, "SET", key, value)
	}
	return err
}

func SetNx(ctx *web.WebContext, key, value string, ttl int64) (interface{}, error) {
	var err error
	var resp interface{}
	if ttl > 0 {
		resp, err = sampleDo(ctx, "SET", key, value, "EX", ttl, "NX")
	} else {
		resp, err = sampleDo(ctx, "SET", key, value, "NX")
	}
	return resp, err
}

func Del(ctx *web.WebContext, key string) error {
	_, err := sampleDo(ctx, "DEL", key)
	if err != nil {
		ctx.Fatalf("Del err=%v", err)
	}
	return err
}

func Incr(ctx *web.WebContext, key string) (int64, error) {
	return sampleDoInt64(ctx, "INCR", key)
}

func Incrby(ctx *web.WebContext, key string, num int) (int64, error) {
	return sampleDoInt64(ctx, "INCRBY", key, num)
}

// Expire
// 当key的值为数字时 返回value为int64,
// 当key的值非数字时，返回value为string
func Expire(ctx *web.WebContext, key string, ttl int) (interface{}, error) {
	var params []interface{}
	var method string

	if ttl > 0 {
		method = "EXPIRE"
		params = append(params, key, ttl)
	} else {
		method = "PERSIST"
		params = append(params, ttl)
	}

	value, err := sampleDo(ctx, method, params...)
	if err != nil {
		return "", err
	}

	return value, nil
}

func Rpush(ctx *web.WebContext, queue, value string) error {
	_, err := sampleDo(ctx, "RPUSH", queue, value)
	return err
}

func Lpop(ctx *web.WebContext, queue string) (string, error) {
	return sampleDoString(ctx, "LPOP", queue)
}

func HSet(ctx *web.WebContext, key, field, value string) error {
	var err error
	_, err = sampleDo(ctx, "HSET", key, field, value)
	return err
}

func HGet(ctx *web.WebContext, key, field string) (string, error) {
	return sampleDoString(ctx, "HGET", key, field)
}

func HDel(ctx *web.WebContext, key string, field string) (string, error) {
	return sampleDoString(ctx, "HDEL", key, field)
}

func SAdd(ctx *web.WebContext, key string, data []string) error {
	var params []interface{}
	params = append(params, key)
	for _, value := range data {
		params = append(params, value)
	}
	_, err := sampleDo(ctx, "SADD", params...)
	return err
}

func SIsMember(ctx *web.WebContext, key string, value string) (int, error) {
	return sampleDoInt(ctx, "SISMEMBER", key, value)
}

func SMembers(ctx *web.WebContext, key string) ([]string, error) {
	return rds.Strings(sampleDo(ctx, "SMEMBERS", key))
}

func SRem(ctx *web.WebContext, key string, data []string) error {
	var params []interface{}
	params = append(params, key)
	for _, value := range data {
		params = append(params, value)
	}
	_, err := sampleDo(ctx, "SREM", params...)
	return err
}
