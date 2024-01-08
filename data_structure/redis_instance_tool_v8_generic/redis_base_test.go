package redis_instance_tool

import (
	"context"
	"fmt"
	"github.com/Xiaoze7713/golib/v3/utils"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestSet(t *testing.T) {
	//rds.Set()
	r := redis.NewClient(&redis.Options{Addr: "39.99.233.6:6379", DB: 7, Password: "redis"})
	p, _ := NewKeyPool[string, string]("ttt", "-", r, false, 0)
	err := p.Set(context.Background(), "a", "bbbb")
	fmt.Println(err)
	//var s string
	val, err := p.GetAndUnmarshal(context.Background(), "a")
	fmt.Println(err)
	fmt.Println(utils.MustJson(val))
}
