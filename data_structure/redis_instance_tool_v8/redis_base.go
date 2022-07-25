package redis_instance_tool

import (
	"git.singularity-ai.com/backend/library/data_structure/message_queue"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

type RedisToolBase struct {
	client redis.Cmdable
	MarshalInterface
	config  *message_queue.MQConfig
	selfKey string
	sep     string
}

type MarshalInterface interface {
	marshalFunction(interface{}) ([]byte, error)
	unmarshalFunction([]byte, interface{}) error
}

var defaultMarshal = MarshalJson{}

type ToolOption interface {
	Apply(r *RedisToolBase)
}

func (m *RedisToolBase) Init(conf *message_queue.MQConfig) {
	m.config = conf
	if m.MarshalInterface == nil {
		json := MarshalJson{}
		json.Apply(m)
	}
}

func (m *RedisToolBase) SelfKey(s string) (key string) {
	return m.selfKey + m.sep + s
}

func (m *RedisToolBase) SelfSep() (sep string) {
	return m.sep
}

type MarshalJson struct{}

func (m *MarshalJson) Apply(r *RedisToolBase) {
	r.MarshalInterface = m
}

func (m *MarshalJson) marshalFunction(i interface{}) ([]byte, error) {
	res, err := jsoniter.Marshal(i)
	return res, err
}

func (m *MarshalJson) unmarshalFunction(bytes []byte, iPtr interface{}) error {
	return jsoniter.Unmarshal(bytes, iPtr)
}
