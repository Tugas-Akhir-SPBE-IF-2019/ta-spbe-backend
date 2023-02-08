package service

type MessageQueue interface {
	Produce(topic string, body []byte) error
}