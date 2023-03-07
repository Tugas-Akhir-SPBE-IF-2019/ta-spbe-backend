package messagequeue

import "github.com/nsqio/go-nsq"

type Config struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Client interface {
	Produce(topic string, body []byte) error
}

type nsqMQ struct {
	producer *nsq.Producer
	consumer *nsq.Consumer
}

func NewMessageQueueNSQ(producer *nsq.Producer, consumer *nsq.Consumer) (Client, error) {
	return &nsqMQ{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (mq *nsqMQ) Produce(topic string, body []byte) error {
	return mq.producer.Publish(topic, body)
}
