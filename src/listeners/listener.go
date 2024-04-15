package listeners

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type listener struct {
	consumer <-chan amqp.Delivery
	service  func(msg amqp.Delivery)
}

func (r *listener) Run() func() {
	return func() {
		for msg := range r.consumer {
			r.service(msg)
		}
	}
}

type IListeners interface {
	Run() func()
}

func InitListeners(Consumer <-chan amqp.Delivery) IListeners {
	return &listener{consumer: Consumer}
}
