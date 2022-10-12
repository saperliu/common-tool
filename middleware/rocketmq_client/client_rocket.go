package rocketmq_client

import (
	"common-tool/common"
	"common-tool/logger"
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
)

type RocketClient struct {
	Address        string
	ProducerGroup  string
	ConsumerGroup  string
	ProducerClient rocketmq.Producer
	ConsumerClient rocketmq.PushConsumer
	Context        context.Context
}

/**
 * 新建rocket客户端
 */
func (rocket *RocketClient) NewRocketProducer(instanceName string) {
	if instanceName == "" {
		//实例名如果相同，则认为是同一个消费者。如果同名且在同一个消费组，则会有一个会不消费信息。
		instanceName = common.RandomString(10)
	}
	server, _ := primitive.NewNamesrvAddr(rocket.Address)
	prod, _ := rocketmq.NewProducer(
		producer.WithNameServer(server),
		producer.WithRetry(3),
		producer.WithDefaultTopicQueueNums(10),
		producer.WithGroupName(rocket.ProducerGroup),
		producer.WithInstanceName(instanceName),
	)
	err := prod.Start()
	if err != nil {
		logger.Error("start producer error: ", err)
		os.Exit(1)
	}
	rocket.ProducerClient = prod
	rocket.Context = context.Background()
}
func (rocket *RocketClient) NewRocketConsumer(instanceName string) {
	if instanceName == "" {
		//实例名如果相同，则认为是同一个消费者。如果同名且在同一个消费组，则会有一个会不消费信息。
		instanceName = common.RandomString(10)
	}
	server, _ := primitive.NewNamesrvAddr(rocket.Address)
	prod, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(rocket.ConsumerGroup),
		consumer.WithNameServer(server),
		consumer.WithAutoCommit(true),
		consumer.WithRetry(3),
		consumer.WithInstance(instanceName),
	)

	//err := rocket.ConsumerClient.Subscribe("test", consumer.MessageSelector{}, func(ctx context.Context,
	//	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	//	for i := range msgs {
	//		logger.Error("subscribe callback: %v \n", msgs[i])
	//	}
	//
	//	return consumer.ConsumeSuccess, nil
	//})
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	rocket.ConsumerClient = prod
	rocket.Context = context.Background()
}

func (rocket *RocketClient) NewConsumer(instanceName string, groupId string) rocketmq.PushConsumer {
	if instanceName == "" {
		//实例名如果相同，则认为是同一个消费者。如果同名且在同一个消费组，则会有一个会不消费信息。
		instanceName = common.RandomString(10)
	}
	if groupId == "" {
		groupId = rocket.ConsumerGroup
	}
	server, _ := primitive.NewNamesrvAddr(rocket.Address)
	prod, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupId),
		consumer.WithNameServer(server),
		consumer.WithAutoCommit(true),
		consumer.WithRetry(3),
		consumer.WithInstance(instanceName),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)

	//err := rocket.ConsumerClient.Subscribe("test", consumer.MessageSelector{}, func(ctx context.Context,
	//	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	//	for i := range msgs {
	//		logger.Error("subscribe callback: %v \n", msgs[i])
	//	}
	//
	//	return consumer.ConsumeSuccess, nil
	//})
	//if err != nil {
	//	logger.Error(err.Error())
	//}

	rocket.Context = context.Background()
	return prod
}

/**
 * 发送消息
 */
func (rocket *RocketClient) SendMessage(topic string, message []byte) error {
	msg := &primitive.Message{
		Topic: topic,
		Body:  message,
	}
	res, err := rocket.ProducerClient.SendSync(rocket.Context, msg)

	if err != nil {
		logger.Error("send message error: %s ", err)
		return err
	} else {
		logger.Info("send message success: result=%s ", res.String())
	}
	return nil
}
func (rocket *RocketClient) Close() {
	err := rocket.ProducerClient.Shutdown()
	if err != nil {
		logger.Error("ProducerClient Shutdown error: %s ", err)
	}
	err = rocket.ConsumerClient.Shutdown()
	if err != nil {
		logger.Error("send message error: %s ", err)
	}
}
