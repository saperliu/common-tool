package kafka_client

import (
	"errors"
	"github.com/childe/healer"
	"github.com/saperliu/common-tool/logger"
	"github.com/segmentio/kafka-go"
	"strings"
)

type ClientKafka struct {
	ProducerTopic string
	ConsumerTopic string
	BrokerList    []string
	GroupId       string
	Producer      *healer.Producer
	Consumer      <-chan *healer.FullMessage
	SegConsumer   *kafka.Reader
}

func (ka *ClientKafka) NewKafkaConsumer(topic string) *kafka.Reader {
	// make a new reader that consumes from topic-A
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  ka.BrokerList,
		GroupID:  ka.GroupId,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	ka.SegConsumer = r
	return r
	//for {
	//	m, err := r.ReadMessage(context.Background())
	//	if err != nil {
	//		break
	//	}
	//	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	//}

	//r.Close()
}
func (ka *ClientKafka) NewConsumer(topic string) (<-chan *healer.FullMessage, error) {
	consumerConfig := healer.DefaultConsumerConfig()
	consumerConfig.BootstrapServers = strings.Join(ka.BrokerList, ",")
	consumerConfig.FromBeginning = false
	consumerConfig.GroupID = ka.GroupId
	consumerConfig.ClientID = "rule-engine"
	var messageChan chan *healer.FullMessage
	messageChan = make(chan *healer.FullMessage, 100)
	consumer, err := healer.NewConsumer(consumerConfig, topic)
	if err != nil {
		logger.Error("could not init GroupConsumer: %s", err)
		return nil, err
	}

	messageChan2, error1 := consumer.Consume(messageChan)
	ka.Consumer = messageChan2
	return messageChan2, error1
}

func (ka *ClientKafka) NewProducer(topic string) (pro *healer.Producer, err error) {

	configMap := make(map[string]interface{})
	configMap["bootstrap.servers"] = strings.Join(ka.BrokerList, ",")
	config, err := healer.GetProducerConfig(configMap)
	if err != nil {
		logger.Error("coult not create producer config: %s", err)
	}
	config.Retries = 3

	producer := healer.NewProducer(topic, config)
	if producer == nil {
		logger.Error("could not create producer")
		return nil, errors.New("could not create producer")
	}
	ka.Producer = producer
	//defer producer.Close()
	return producer, nil
}

func (ka *ClientKafka) SendMessage(key string, value string) error {
	err := ka.Producer.AddMessage([]byte(key), []byte(value))
	if err != nil {
		logger.Error("send message error  %v ", err)
		return err
	}
	return nil
}
