package ali_rocketmq

import (
	"common-tool/common"
	"common-tool/logger"
	"context"
	"github.com/aliyunmq/mq-http-go-sdk"
)

type AliRocketClient struct {
	ProducerGroup         string
	ConsumerGroup         string
	Endpoint              string
	AccessKey             string
	SecretKey             string
	ProducerClient        mq_http_sdk.MQProducer
	ProducerRealClient    mq_http_sdk.MQProducer
	ProducerHistoryClient mq_http_sdk.MQProducer
	ConsumerClient        mq_http_sdk.MQConsumer
	Context               context.Context
}

// 设置HTTP接入域名（此处以公共云生产环境为例）
//endpoint := "${HTTP_ENDPOINT}"
// AccessKey 阿里云身份验证，在阿里云服务器管理控制台创建
//accessKey := "${ACCESS_KEY}"
// SecretKey 阿里云身份验证，在阿里云服务器管理控制台创建
//secretKey := "${SECRET_KEY}"
// 所属的 Topic
//topic := "${TOPIC}"
// Topic所属实例ID，默认实例为空
//instanceId := "${INSTANCE_ID}"

/**
 * 新建rocket客户端
 */
func (rocket *AliRocketClient) NewRocketProducer(instanceName string, topic string) {
	client := mq_http_sdk.NewAliyunMQClient(rocket.Endpoint, rocket.AccessKey, rocket.SecretKey, "")
	if instanceName == "" {
		//实例名如果相同，则认为是同一个消费者。如果同名且在同一个消费组，则会有一个会不消费信息。
		instanceName = common.RandomString(10)
	}
	mqProducer := client.GetProducer(instanceName, topic)
	rocket.ProducerClient = mqProducer
	rocket.Context = context.Background()
}

/**
 * 新建rocket消费客户端
 */
func (rocket *AliRocketClient) NewRocketConsume(instanceName string, topic string) mq_http_sdk.MQConsumer {
	client := mq_http_sdk.NewAliyunMQClient(rocket.Endpoint, rocket.AccessKey, rocket.SecretKey, "")
	if instanceName == "" {
		//实例名如果相同，则认为是同一个消费者。如果同名且在同一个消费组，则会有一个会不消费信息。
		instanceName = common.RandomString(10)
	}
	mqConsumer := client.GetConsumer(instanceName, topic, rocket.ConsumerGroup, "")
	rocket.ConsumerClient = mqConsumer
	return mqConsumer
}

/**
 * 发送消息
 */
func (rocket *AliRocketClient) SendMessage(message string, messageTag string, key string) {
	var msg mq_http_sdk.PublishMessageRequest
	msg = mq_http_sdk.PublishMessageRequest{
		MessageBody: message,             //消息内容
		MessageTag:  messageTag,          // 消息标签
		Properties:  map[string]string{}, // 消息属性
	}
	// 设置KEY
	msg.MessageKey = key
	// 设置属性
	//msg.Properties["a"] = strconv.Itoa(i)
	ret, err := rocket.ProducerClient.PublishMessage(msg)
	if err != nil {
		logger.Error("Publish Message error: ", err)
	} else {
		logger.Info("Publish MessageId:%s, BodyMD5:%s, ", ret.MessageId, ret.MessageBodyMD5)
	}
}

/**
 * 发送消息
 */
func (rocket *AliRocketClient) SendRealMessage(message string, messageTag string, key string) {
	var msg mq_http_sdk.PublishMessageRequest
	msg = mq_http_sdk.PublishMessageRequest{
		MessageBody: message,             //消息内容
		MessageTag:  messageTag,          // 消息标签
		Properties:  map[string]string{}, // 消息属性
	}
	// 设置KEY
	msg.MessageKey = key
	// 设置属性
	//msg.Properties["a"] = strconv.Itoa(i)
	ret, err := rocket.ProducerRealClient.PublishMessage(msg)
	if err != nil {
		logger.Error("Publish Message error: ", err)
	} else {
		logger.Info("Publish MessageId:%s, BodyMD5:%s, ", ret.MessageId, ret.MessageBodyMD5)
	}
}

/**
 * 发送消息
 */
func (rocket *AliRocketClient) SendHistoryMessage(message string, messageTag string, key string) {
	var msg mq_http_sdk.PublishMessageRequest
	msg = mq_http_sdk.PublishMessageRequest{
		MessageBody: message,             //消息内容
		MessageTag:  messageTag,          // 消息标签
		Properties:  map[string]string{}, // 消息属性
	}
	// 设置KEY
	msg.MessageKey = key
	// 设置属性
	//msg.Properties["a"] = strconv.Itoa(i)
	ret, err := rocket.ProducerHistoryClient.PublishMessage(msg)
	if err != nil {
		logger.Error("Publish Message error: ", err)
	} else {
		logger.Info("Publish MessageId:%s, BodyMD5:%s, ", ret.MessageId, ret.MessageBodyMD5)
	}
}
