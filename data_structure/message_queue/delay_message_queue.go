package message_queue

import (
	"context"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-stack/stack"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type DMQInstance interface {
	Add(ctx context.Context, key string, valueIf interface{}, delayDur int64, expireDur int64) (err error)
	Count(ctx context.Context, timeMsStart, timeMsEnd int64) (count int, err error)
	PopL(ctx context.Context, timeMsStart, timeMsEnd int64) (resList []string, err error)
	Name() string
	Config() *MQConfig
}

type DqMsg struct {
	Key       string
	Value     interface{}
	DelayDur  int64
	ExpireDur int64
}

type DMQManager struct {
	instance        DMQInstance
	handlerFunction MQMsgHandlerFunc
	consumerCtx     context.Context
	consumerCancel  context.CancelFunc
	producerCtx     context.Context
	producerCancel  context.CancelFunc
	ctx             context.Context
	//msgIn           chan interface{}
	msgOut chan *DqMsg
}

func NewDMQManager(instance DMQInstance, handler MQMsgHandlerFunc) (mqMgr *DMQManager, err error) {
	if handler == nil {
		handler = DoNoting
	}
	mqMgr = &DMQManager{
		instance: instance,
		//serializationFunc:   serializationFunc,
		//deserializationFunc: deserializationFunc,
		handlerFunction: handler,
		ctx:             context.Background(),
		//msgIn:           make(chan interface{}, 10),
		msgOut: make(chan *DqMsg, 10),
	}
	mqMgr.consumerCtx, mqMgr.consumerCancel = context.WithCancel(mqMgr.ctx)
	mqMgr.producerCtx, mqMgr.producerCancel = context.WithCancel(mqMgr.ctx)
	//mqMgr.Start()
	return mqMgr, nil
}

func (m *DMQManager) LoopConsumer() {
	defer func() {
		if err := recover(); err != any(nil) {
			xlog.Errorf("run mq message err %v", err)
		}
	}()
	xlog.Infof("start run mq consumer %v", m.instance.Name())
	for {
		select {
		//case msg := <-m.msgIn:
		case <-m.consumerCtx.Done():
			xlog.Debugf("consumer close")
			return
		default:
			t := time.Now().UnixMilli()
			msgList, err := m.instance.PopL(m.consumerCtx, 0, t)
			if err != nil {
				xlog.Errorf("%v read msg error", err)
				//m.ctx.Done()
				//return
				time.Sleep(time.Millisecond * 500)
				continue
			}
			if len(msgList) <= 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			for _, msgIF := range msgList {
				msg := msgIF
				go func(newMsg interface{}) {
					// 开始span
					defer func() {
						if r := recover(); r != any(nil) {
							xlog.Errorf("panic msg=%v,  recover=%v %v", stack.Trace())
						}
					}()
					m.handlerFunction(newMsg)
				}(msg)
				//m.msgIn <- msg
			}
		}
	}
}
func (m *DMQManager) LoopProducer() {
	for {
		select {
		case msg := <-m.msgOut:
			err := m.instance.Add(m.producerCtx, msg.Key, msg.Value, msg.DelayDur, msg.ExpireDur)
			if err != nil {
				msgStr, err1 := jsoniter.MarshalToString(msg)
				if err1 != nil {
					xlog.Errorf("send failed known msg type, err=%v", err1)
				}
				xlog.Errorf("send failed msg=%v", msgStr)
			}
			continue
		case <-m.producerCtx.Done():
			xlog.Debugf(" close")
		}
	}
}
func (m *DMQManager) Send(msg *DqMsg) error {
	m.msgOut <- msg
	return nil
}

func (m *DMQManager) Close() {
	m.producerCancel()
	m.consumerCancel()
}

//func (m *DMQManager) Run() {
//	defer func() {
//		if err := recover(); err != any(nil) {
//			xlog.Errorf("run mq message err %v", err)
//		}
//		go m.LoopConsumer()
//		go m.LoopProducer()
//		select {
//		case <-m.ctx.Done():
//			m.Close()
//		}
//	}()
//}
//

func (m *DMQManager) RunProducer() {
	defer func() {
		if err := recover(); err != any(nil) {
			m.ctx.Done()
			xlog.Errorf("run mq message err %v", err)
		}
		xlog.Debugf("start run mq producer %v", m.instance.Name())
		go m.LoopProducer()
		select {
		case <-m.ctx.Done():
			m.Close()
		}
	}()
}

func (m *DMQManager) RunConsumer() {
	defer func() {
		go m.LoopConsumer()
		select {
		case <-m.ctx.Done():
			m.Close()
		}
	}()
}
