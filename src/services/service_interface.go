package services

import amqp "github.com/rabbitmq/amqp091-go"

type IService[P any] interface {
	GetProcessMsg() func(msg amqp.Delivery)
}
