package kafka

import (
	"context"
	"errors"
	"git.singularity-ai.com/backend/library/v2/utils"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"sync"
)

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
				if data != nil && handleFunc != nil {
					handleFunc(data)
				} else {
					xlog.Errorf("kafka[%v] error msg %v", utils.MustJson(m.conf), utils.MustJson(data))
				}
			default:
				xlog.Warnf("kafka[%v] unknown msg %v", utils.MustJson(m.conf), utils.MustJson(e))
			}
		case err = <-m.ErrChan:
			xlog.Error(err)
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
			err = m.p.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &m.conf.Topic, Partition: kafka.PartitionAny}, Value: value}, nil)
			if err != nil {
				xlog.Error(err)
				//m.ErrChan <- err
			}
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

	p, err := NewKafkaProducer(config)
	if err != nil {
		xlog.Errorf("NewProducer failed, err:%v", err)
		return nil, err
	}
	producer.p = p
	return producer, nil
}

func (p *Producer) SendMsg(msgByte []byte) error {
	p.ProduceChan <- msgByte
	return nil
}

func NewKafkaProducer(cfg *Config) (p *kafka.Producer, err error) {
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request": "true",
		"message.max.bytes":   1000000,
		"linger.ms":           10,
		"retries":             30,
		"retry.backoff.ms":    1000,
		"acks":                "1"}
	kafkaconf.SetKey("bootstrap.servers", cfg.Broker)

	switch cfg.SecurityProtocol {
	case "plaintext":
		kafkaconf.SetKey("security.protocol", "plaintext")
	case "sasl_ssl":
		kafkaconf.SetKey("security.protocol", "sasl_ssl")
		kafkaconf.SetKey("ssl.ca.location", cfg.SslCaLocation)
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
	case "sasl_plaintext":
		kafkaconf.SetKey("sasl.mechanism", "PLAIN")
		kafkaconf.SetKey("security.protocol", "sasl_plaintext")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
	default:
		err = kafka.NewError(kafka.ErrUnknownProtocol, "unknown protocol", true)
		return nil, err
	}

	producer, err := kafka.NewProducer(kafkaconf)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
