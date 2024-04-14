package rabbits

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func InitQueue(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func Publish[T any](ch *amqp.Channel) func(queueName string, msg T) error {
	return func(queueName string, msg T) error {
		body, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = ch.PublishWithContext(
			ctx,
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		return err
	}
}

func Consume(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
