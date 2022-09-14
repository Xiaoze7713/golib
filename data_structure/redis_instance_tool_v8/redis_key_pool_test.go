package redis_instance_tool

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestTokenPool(t *testing.T) {
	rdsPool := &redis.Pool{
		MaxIdle:   20,
		MaxActive: 100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "39.99.156.201:6379")
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", "redis"); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", 8); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	rdsPool.Close()
	//kpl := TokenPool{KeyPool{RedisCli: RedisCli}}
	//err := kpl.Set("123", "123", 30, true)
	//fmt.Printf("%v\n", err)
	//res, err := kpl.Get("123")
	//fmt.Printf("%v %v\n", res, err)
}

func TestCtx(t *testing.T) {
	bc_ctx := context.Background()
	fmt.Printf("%+v\n", &bc_ctx)
	valueCtx := context.WithValue(context.Background(), "hi", "mingze")
	fmt.Printf("%+v\n", &valueCtx)
	fmt.Printf("%+v\n", valueCtx.Value("hi"))
	newVCtx, cancel := context.WithCancel(valueCtx)
	fmt.Printf("%+v\n", newVCtx.Value("hi"))
	fmt.Printf("%+v\n", context.Background().Value("hi"))

	defer cancel()
}

func TestJobKey(t *testing.T) {

}
