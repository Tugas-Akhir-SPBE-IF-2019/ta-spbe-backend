package messagequeue

import "github.com/nsqio/go-nsq"

type NSQ struct {
	Producer *nsq.Producer
	Consumer *nsq.Consumer
}

func (mq *NSQ) Produce(topic string, body []byte) error {
	return mq.Producer.Publish(topic, body)
}
