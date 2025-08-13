package redis_instance_tool

import (
	"errors"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

var (
	WarnKeysCountUse = errors.New("warn empty key, can`t use keys")
)

type RedisToolBase struct {
	client redis.Cmdable
	MarshalInterface
	//config  *message_queue.MQConfig
	selfKey string
	sep     string
}

type MarshalInterface interface {
	marshalFunction(interface{}) (string, error)
	unmarshalFunction(string, interface{}) error
}

var defaultMarshal = &MarshalJson{}

type ToolOption interface {
	Apply(r *RedisToolBase)
}

func (m *RedisToolBase) Init() {
	//m.config = conf
	if m.MarshalInterface == nil {
		json := MarshalJson{}
		json.Apply(m)
	}
}

func (m *RedisToolBase) SelfKey(s string) (key string) {
	return m.Prefix() + s
}

func (m *RedisToolBase) Prefix() string {
	return m.selfKey + m.sep
}
func (m *RedisToolBase) SelfSep() (sep string) {
	return m.sep
}

type MarshalJson struct{}

func (m *MarshalJson) Apply(r *RedisToolBase) {
	r.MarshalInterface = m
}

func (m *MarshalJson) marshalFunction(i interface{}) (string, error) {
	res, err := jsoniter.MarshalToString(i)
	return res, err
}

func (m *MarshalJson) unmarshalFunction(s string, iPtr interface{}) error {
	return jsoniter.UnmarshalFromString(s, iPtr)
}
