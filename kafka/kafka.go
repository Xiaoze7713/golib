/**
 * @Author: wenliangzhang
 * @Description:
 * @File: kafak
 * @Version: 1.0.0
 * @Date: 2022/4/19 7:59 PM
 */
package kafka

import (
	"context"
	"fmt"
	"git.singularity-ai.com/backend/library/log"
	"github.com/BurntSushi/toml"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"io/ioutil"
	"os"
	"sync"
)

type KafkaConfig struct {
	Name   string `toml:"name"`
	Broker string `toml:"broker"`
	Topic  string `toml:"topic"`
}

type Producer struct {
	ProduceChan chan []byte
	ErrChan     chan error
	Closer      context.CancelFunc
}

type Consumer struct {
	ConsumeChan chan *kafka.Message
	ErrChan     chan error
	Closer      context.CancelFunc
}

var defaultKafkaConfigPath = "conf/service/kafka.toml"

var configs map[string]KafkaConfig

var producers sync.Map
var consumers sync.Map

// Init
//	配置文件格式示例：
//	[receive_ws]
//	name = "receive_ws"
//	broker = "172.21.77.44:19092"
//	topic = "receive_ws_msg"
func Init(filePath string) {
	var config map[string]KafkaConfig
	if filePath == "" {
		filePath = defaultKafkaConfigPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Infoln("kafka.toml not exist")
		return
	}
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("read kafka.toml failed,err=%s", err.Error())
		return
	}
	var confStrRaw = string(bs)
	if _, err := toml.Decode(confStrRaw, &config); err != nil {
		panic(fmt.Sprintf("Can't load config file, %s", err.Error()))
	}

	if len(config) > 0 {
		configs = make(map[string]KafkaConfig)
		for name, item := range config {
			configs[name] = item
		}
	}
}

func GetConfig(name string) KafkaConfig {
	return configs[name]
}

func GetProducer(name string) Producer {
	producer, ok := producers.Load(name)
	if !ok {
		config := GetConfig(name)
		producer, err := NewProducer(config.Broker, config.Topic)
		if err != nil {
			return Producer{}
		} else {
			producers.Store(name, producer)
			return producer
		}
	}
	return producer.(Producer)
}

func NewProducer(broker string, topic string) (Producer, error) {
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup

	var producer Producer
	producer.ProduceChan = make(chan []byte, 100)
	producer.ErrChan = make(chan error, 100)
	ctx, closer := context.WithCancel(context.Background())
	producer.Closer = closer

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Errorf("NewProducer failed, err:%v", err)
		return Producer{}, err
	}

	wg1.Add(1)
	go func() {
		for value := range producer.ProduceChan {
			p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}
		}
		wg1.Done()
	}()

	wg2.Add(1)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					producer.ErrChan <- m.TopicPartition.Error
				}
			default:
			}
		}
		wg2.Done()
	}()

	go func() {
		select {
		case <-ctx.Done():
			// order matters
			close(producer.ProduceChan)
			wg1.Wait()
			p.Close()
			wg2.Wait()
			close(producer.ErrChan)
		}
	}()

	return producer, nil
}

func (p *Producer) SendMsg(msgByte []byte) error {
	p.ProduceChan <- msgByte
	return nil
}

func GetConsumer(name string) Consumer {
	consumer, ok := consumers.Load(name)
	if !ok {
		config := GetConfig(name)
		consumer, err := NewConsumer(config.Broker, config.Topic, name)
		if err != nil {
			return Consumer{}
		} else {
			producers.Store(name, consumer)
			return consumer
		}
	}
	return consumer.(Consumer)
}

func NewConsumer(broker, topic, group string) (Consumer, error) {
	var wg1 sync.WaitGroup

	// 如果group为空则读不到
	if group == "" {
		group = "group"
	}

	var consumer Consumer
	consumer.ConsumeChan = make(chan *kafka.Message, 100)
	consumer.ErrChan = make(chan error, 100)
	ctx, closer := context.WithCancel(context.Background())
	consumer.Closer = closer

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		//"auto.offset.reset":               "earliest",
	})

	if err != nil {
		log.Errorf("NewConsumer failed, err:%v", err)
		return Consumer{}, err
	}

	err = c.Subscribe(topic, nil)
	if err != nil {
		log.Errorf("Subscribe failed, err:%v", err)
		return Consumer{}, err
	}

	wg1.Add(1)
	go func() {
		for ev := range c.Events() {
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				c.Unassign()
			case *kafka.Message:
				consumer.ConsumeChan <- e
			case kafka.Error:
				consumer.ErrChan <- e
			case kafka.OffsetsCommitted:
			default:
				log.Printf("unknown event %+v", e)
			}
		}
		wg1.Done()
	}()

	go func() {
		select {
		case <-ctx.Done():
			c.Close()
			wg1.Wait()
			close(consumer.ConsumeChan)
			close(consumer.ErrChan)
		}
	}()

	return consumer, nil
}
