package redis_instance_tool

import (
	"git.singularity-ai.com/backend/library/data_structure/message_queue"
	"github.com/go-redis/redis/v8"
)

type RedisToolBase struct {
	client  redis.Cmdable
	config  *message_queue.MQConfig
	selfKey string
	sep     string
}

func (m *RedisToolBase) Init(conf *message_queue.MQConfig) {
	m.config = conf
}

func (m *RedisToolBase) SelfKey(s string) (key string) {
	return m.selfKey + m.sep + s
}

func (m *RedisToolBase) SelfSep() (sep string) {
	return m.sep
}

func (m *RedisToolBase) Marshal(i interface{}) (string, error) {
	return "", nil
}

func (m *RedisToolBase) UnMarshal(s string, iPtr interface{}) error {
	return nil
}
