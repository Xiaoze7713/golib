package kafka

import (
	"context"
	"errors"
	"git.singularity-ai.com/backend/library/v2/utils"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

type Consumer struct {
	ConsumeChan chan *kafka.Message
	ErrChan     chan error
	conf        *Config
	ctxCancel   context.CancelFunc
	ctx         context.Context
	c           *kafka.Consumer
	w           *sync.WaitGroup
}

func (m *Consumer) eventListen(handleFunc func(message *kafka.Message)) (err error) {
	err = m.c.Subscribe(m.conf.Topic, nil)
	if err != nil {
		xlog.Errorf("Subscribe failed, err:%v", err)
		return err
	}
	for {
		select {
		case <-m.ctx.Done():
			if errors.Is(m.ctx.Err(), context.Canceled) {
				xlog.Infof("kafka[%v] event listen has cancel!", utils.MustJson(m.conf))
			} else {
				xlog.Infof("kafka[%v] event listen has timeout!", utils.MustJson(m.conf))
			}
			return
		case e := <-m.c.Events():
			switch ev := e.(type) {
			case kafka.AssignedPartitions:
				err := m.c.Assign(ev.Partitions)
				if err != nil {
					xlog.Infof("kafka[%v] event listen err %v", utils.MustJson(m.conf), err)
					return
				}
			case kafka.RevokedPartitions:
				err := m.c.Unassign()
				if err != nil {
					xlog.Infof("kafka[%v] event listen err %v", utils.MustJson(m.conf), err)
					return
				}
			case *kafka.Message:
				m.ConsumeChan <- ev
			case kafka.Error:
				m.ErrChan <- ev
			case kafka.OffsetsCommitted:
			default:
				xlog.Warnf("kafka[%v] event listen, unknown event %+v", e)
			}
		}
	}
}

func (m *Consumer) Run(errHandleFunc func(message *kafka.Message), msgHandler func(message *kafka.Message)) (err error) {
	m.w.Add(1)
	go func() {
		xlog.Infof("kafka[%v] event listen running", utils.MustJson(m.conf))
		err = m.eventListen(errHandleFunc)
		if err != nil {
			xlog.Infof("kafka[%v] event listen running", utils.MustJson(m.conf))
			return
		}
		m.w.Done()
	}()
	return
}

func NewConsumer(ctx context.Context, config *Config) (consumer *Consumer, err error) {
	// 如果group为空则默认
	if config.Group == "" {
		config.Group = "default"
		xlog.Warnf("consumer group use %v", config.Group)
	}
	consumer.ConsumeChan = make(chan *kafka.Message, 1024*8)
	consumer.ErrChan = make(chan error, 1024*8)
	consumer.ctx, consumer.ctxCancel = context.WithCancel(context.Background())
	consumer.w = &sync.WaitGroup{}
	consumer.conf = config

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               consumer.conf.Broker,
		"group.id":                        consumer.conf.Group,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		//"auto.offset.reset":               "earliest",
	})

	if err != nil {
		xlog.Errorf("NewConsumer failed, err:%v", err)
		return nil, err
	}

	consumer.c = c
	return consumer, nil
}
