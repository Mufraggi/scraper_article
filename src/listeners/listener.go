package listeners

import (
	"fmt"
	"github.com/Mufraggi/scraper_article/src/mongo_repo"
	"github.com/Mufraggi/scraper_article/src/rabbits"
	amqp "github.com/rabbitmq/amqp091-go"
)

type listener[R any, P any, B any] struct {
	mongoRepo mongo_repo.IRepository[R]
	publish   func(queueName string, msg P) error
	scraper   func(url string) B
	Consumer  <-chan amqp.Delivery
}

func (r *listener[R, P, B]) Run() func() {
	return func() {
		for msg := range r.Consumer {
			fmt.Printf(msg)
		}
	}
}

type IListeners[R any, P any, B any] interface {
	Run() func()
}

func InitListeners[R any, P any, B any](ch *amqp.Channel, mongoRepo mongo_repo.IRepository[R],
	scraper func(url string) B,
	Consumer <-chan amqp.Delivery) IListeners[R, P, B] {
	publish := rabbits.Publish[P](ch)
	return &listener[R, P, B]{mongoRepo: mongoRepo, publish: publish, scraper: scraper, Consumer: Consumer}
}
