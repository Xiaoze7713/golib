package message_queue

import (
	"context"
	"fmt"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"github.com/go-stack/stack"
	jsoniter "github.com/json-iterator/go"
)

type MQConfig struct {
	Topic        string `json:"topic" toml:"topic"`
	InstanceName string `json:"instance_name" toml:"instance_name"`
	InstanceType string `json:"instance_type" toml:"instance_type"`
	Broker       string `json:"broker" toml:"broker"`
	Group        string `json:"group" toml:"group"`
}

func (m *MQConfig) String() string {
	s, _ := jsoniter.MarshalToString(m)
	return fmt.Sprintf("%s", s)
}

type MQInstance[T any] interface {
	Push(ctx context.Context, data T) error
	Pop(ctx context.Context, block bool) (T, error)
	Name() string
	Config() *MQConfig
}

type MQMsgHandlerFunc func(interface{})

func DoNoting(i interface{}) {}

type MQManager[T any] struct {
	instance        MQInstance[T]
	handlerFunction MQMsgHandlerFunc
	consumerCtx     context.Context
	consumerCancel  context.CancelFunc
	producerCtx     context.Context
	producerCancel  context.CancelFunc
	ctx             context.Context
	msgIn           chan T
	msgOut          chan T
}

func NewMQManager[T any](instance MQInstance[T], handler MQMsgHandlerFunc) (mqMgr *MQManager[T], err error) {
	mqMgr = &MQManager[T]{
		instance:        instance,
		handlerFunction: handler,
		ctx:             context.Background(),
		msgIn:           make(chan T, 1000),
		msgOut:          make(chan T, 10),
	}
	mqMgr.consumerCtx, mqMgr.consumerCancel = context.WithCancel(mqMgr.ctx)
	mqMgr.producerCtx, mqMgr.producerCancel = context.WithCancel(mqMgr.ctx)
	//mqMgr.Start()
	return mqMgr, nil
}

func (m *MQManager[T]) LoopConsumer() {
	xlog.Infof("start consumer topic=%v, type=%v, broker=%v, ", m.instance.Config().Topic, m.instance.Config().InstanceType, m.instance.Config().Broker)
	for {
		select {
		case msg := <-m.msgIn:
			go func(newMsg interface{}) {
				// 开始span
				defer func() {
					if r := recover(); r != any(nil) {
						xlog.Errorf("panic msg=%v,  recover=%v %v", msg, r, stack.Trace())
					}
				}()
				m.handlerFunction(newMsg)
			}(msg)
		case <-m.consumerCtx.Done():
			xlog.Infof("consumer close")
		default:
			msgI, err := m.instance.Pop(m.consumerCtx, true)
			if err != nil {
				xlog.Errorf("%v read msg error", err)
				//m.ctx.Done()
				continue
			}
			//msg := msgI.(string)
			m.msgIn <- msgI
		}
	}
}
func (m *MQManager[T]) LoopProducer() {
	for {
		select {
		case msg := <-m.msgOut:
			err := m.instance.Push(m.producerCtx, msg)
			if err != nil {
				msgStr, err1 := jsoniter.MarshalToString(msg)
				if err1 != nil {
					xlog.Debugf("%v %v", msgStr, err1)
				}
				xlog.Errorf("push %v error %v", msgStr, err)
			}
			continue
		case <-m.producerCtx.Done():
			xlog.Infof("%v close", m.instance.Config().InstanceName)
		}
	}
}

func (m *MQManager[T]) Send(msg T) error {
	m.msgOut <- msg
	return nil
}

func (m *MQManager[T]) Close() {
	m.producerCancel()
	m.consumerCancel()
}

//func (m *MQManager) Run() {
//	defer func() {
//		go m.LoopConsumer()
//		go m.LoopProducer()
//		select {
//		case <-m.ctx.Done():
//			m.Close()
//		}
//	}()
//}
//

func (m *MQManager[T]) RunProducer() {
	go func() {
		if r := recover(); r != any(nil) {
			m.ctx.Done()
			xlog.Errorf("run mq producer %v, %v %v", m.instance.Config().Topic, r, stack.Trace())
		}
		m.LoopProducer()
	}()
	select {
	case <-m.ctx.Done():
		m.Close()
	}
}

func (m *MQManager[T]) RunConsumer() {
	go func() {
		if r := recover(); r != any(nil) {
			m.ctx.Done()
			xlog.Errorf("run mq consumer %v %v %v", m.instance.Config().Topic, r, stack.Trace())
		}
		m.LoopConsumer()
	}()
	select {
	case <-m.ctx.Done():
		m.Close()
	}
}
