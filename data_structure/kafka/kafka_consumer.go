package kafka

import (
	"context"
	"errors"
	"git.singularity-ai.com/backend/library/v2/log"
	"git.singularity-ai.com/backend/library/v2/utils"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

type Config struct {
	Name   string `toml:"name" json:"name"`
	Broker string `toml:"broker" json:"broker"`
	Topic  string `toml:"topic" json:"topic"`
	Group  string `toml:"group" json:"group"`
}

type Producer struct {
	ProduceChan chan []byte
	ErrChan     chan error
	conf        *Config
	ctxCancel   context.CancelFunc
	ctx         context.Context
	p           *kafka.Producer
	w           *sync.WaitGroup
}

func (m *Producer) eventListen(handleFunc func(message *kafka.Message)) (err error) {
	for {
		select {
		case <-m.ctx.Done():
			if errors.Is(m.ctx.Err(), context.Canceled) {
				xlog.Infof("kafka[%v] event listen has cancel!", utils.MustJson(m.conf))
			} else {
				xlog.Infof("kafka[%v] event listen has timeout!", utils.MustJson(m.conf))
			}
			return
		case e := <-m.p.Events():
			switch ev := e.(type) {
			case *kafka.Message:
				data := ev
				if data != nil {
					handleFunc(data)
				} else {
					xlog.Errorf("kafka[%v] error msg %v", utils.MustJson(m.conf), utils.MustJson(data))
				}
			default:
				xlog.Warnf("kafka[%v] unknown msg %v", utils.MustJson(m.conf), utils.MustJson(e))
			}
		}
	}
}

func (m *Producer) msgLoopProduct() (err error) {
	for {
		select {
		case <-m.ctx.Done():
			if errors.Is(m.ctx.Err(), context.Canceled) {
				xlog.Infof("kafka[%v] ctx has cancel", utils.MustJson(m.conf))
			} else {
				xlog.Infof("kafka[%v] ctx has timeout", utils.MustJson(m.conf))
			}
			return
		case value := <-m.ProduceChan:
			m.p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &m.conf.Topic, Partition: kafka.PartitionAny}, Value: []byte(value)}
		}
	}
}

func (m *Producer) Run(errHandleFunc func(message *kafka.Message)) (err error) {
	m.w.Add(1)
	go func() {
		xlog.Infof("kafka[%v] event listen running", utils.MustJson(m.conf))
		err = m.eventListen(errHandleFunc)
		if err != nil {
			xlog.Errorf("kafka[%v] event listen error %v", utils.MustJson(m.conf), err)
			return
		}
		m.w.Done()
	}()
	m.w.Add(1)
	go func() {
		xlog.Infof("kafka[%v] msg listen running", utils.MustJson(m.conf))
		err = m.msgLoopProduct()
		if err != nil {
			xlog.Errorf("kafka[%v] msg listen error %v", utils.MustJson(m.conf), err)
			return
		}
		m.w.Done()
	}()
	return err
}

func (m *Producer) Close() {
	m.ctxCancel()
	m.w.Wait()
	close(m.ProduceChan)
	close(m.ErrChan)
	m.p.Close()
}

// Init
//	配置文件格式示例：
//	[event]
//	name = "event"
//	broker = "172.21.77.44:19092"
//	topic = "event_service"

func NewProducer(ctx context.Context, config *Config) (producer *Producer, err error) {
	producer = &Producer{}
	producer.ProduceChan = make(chan []byte, 1024*8)
	producer.ErrChan = make(chan error, 1024*8)
	producer.ctx, producer.ctxCancel = context.WithCancel(context.Background())
	producer.w = &sync.WaitGroup{}
	producer.conf = config

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": producer.conf.Broker})
	if err != nil {
		log.Errorf("NewProducer failed, err:%v", err)
		return nil, err
	}
	producer.p = p
	return producer, nil
}

func (p *Producer) SendMsg(msgByte []byte) error {
	p.ProduceChan <- msgByte
	return nil
}
