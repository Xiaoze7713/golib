package redis_instance_tool

import (
	"git.singularity-ai.com/backend/library/data_structure/message_queue"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

type RedisToolBase struct {
	client            redis.Cmdable
	marshalFunction   func(interface{}) ([]byte, error)
	unmarshalFunction func([]byte, interface{}) error
	config            *message_queue.MQConfig
	selfKey           string
	sep               string
}

type ToolOption interface {
	Apply(r *RedisToolBase)
}

func (m *RedisToolBase) Init(conf *message_queue.MQConfig) {
	m.config = conf
	if m.marshalFunction == nil {
		MarshalJson{}.Apply(m)
	}
}

func (m *RedisToolBase) SelfKey(s string) (key string) {
	return m.selfKey + m.sep + s
}

func (m *RedisToolBase) SelfSep() (sep string) {
	return m.sep
}

type MarshalJson struct{}

func (m MarshalJson) Apply(r *RedisToolBase) {
	r.marshalFunction = jsoniter.Marshal
	r.unmarshalFunction = jsoniter.Unmarshal
}

func (m *RedisToolBase) Marshal(i interface{}) (string, error) {
	bs, err := m.marshalFunction(i)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (m *RedisToolBase) UnMarshal(s string, iPtr interface{}) error {
	return m.unmarshalFunction([]byte(s), iPtr)
}
