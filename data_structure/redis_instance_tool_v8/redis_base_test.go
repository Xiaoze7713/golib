package redis_instance_tool

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestSet(t *testing.T) {
	//rds.Set()
	r := redis.NewClient(&redis.Options{Addr: "39.99.233.6:6379", DB: 7, Password: "redis"})
	p, _ := NewKeyPool("ttt", "-", r, false, 0)
	err := p.Set(context.Background(), "a", "bbbb")
	fmt.Println(err)
	var s string
	err = p.GetAndUnmarshal(context.Background(), "a", &s)
	fmt.Println(err)
	fmt.Println(s)
}
