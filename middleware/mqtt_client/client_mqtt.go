package mqtt_client

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/saperliu/common-tool/logger"
)

type MqttClient struct {
	ClientId  string
	BrokerUrl string
	Username  string
	Password  string
	//SkipTLSVerification bool
	//NumberOfMessages    int
	//Timeout             time.Duration
	//Retained            bool
	//PublisherQoS        byte
	//SubscriberQoS       byte
	Client *MQTT.Client
}

// define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	logger.Info("TOPIC: %s  MSG: %s", msg.Topic(), msg.Payload())
}

/**
 * 新建mqtt客户端
 */
func (mqttClient *MqttClient) NewMqttClient() *MQTT.Client {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(mqttClient.BrokerUrl)
	opts.SetClientID(mqttClient.ClientId)
	opts.SetUsername(mqttClient.Username)
	opts.SetPassword(mqttClient.Password)
	opts.SetDefaultPublishHandler(f)
	opts.SetAutoReconnect(true)
	//opts.SetConnectRetry(true)
	//opts.SetConnectRetryInterval(1 * time.Second)
	opts.SetConnectionLostHandler(func(client MQTT.Client, msg error) {
		// 如果 playload是json则会报错
		logger.Error("ConnectionLostHandler  TOPIC: %s\n", msg.Error())
	})
	//create and start a client using the above ClientOptions
	clientMq := MQTT.NewClient(opts)
	if token := clientMq.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	mqttClient.Client = &clientMq
	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	//token := clientMq.Subscribe(topicName, mqttClient.SubscriberQoS, func(client MQTT.Client, msg MQTT.Message) {
	//	//fmt.Printf("sub  TOPIC: %s\n", msg.Topic())
	//	logger.Info("MSG: %s\n", msg.Payload())
	//	go mqttClient.Subscrib(msg.Topic(), msg.Payload())
	//
	//})
	//if	token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	//os.Exit(1)
	//}
	return &clientMq
}
