package services

import (
	"encoding/json"
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/Mufraggi/scraper_article/src/mongo_repo"
	"github.com/Mufraggi/scraper_article/src/rabbits"
	amqp "github.com/rabbitmq/amqp091-go"
)

type serviceDetail[P any, R domain.Detail] struct {
	publish   func(queueName string, msg P) error
	scraper   func(url string) R
	mongoRepo mongo_repo.IRepository[R]
}

type IService[P any, R any] interface {
	GetProcessMsg() func(msg amqp.Delivery)
}

func (s *serviceDetail[P, R]) GetProcessMsg() func(msg amqp.Delivery) {
	return func(msg amqp.Delivery) {
		var message domain.DetailRabbitMsg
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			fmt.Printf("error deserialization")
		}
		fmt.Println(message)
		res := s.scraper(message.Url)
		id, err := s.mongoRepo.Create(res)
		if err != nil {
			fmt.Printf("error creating mongo repo")
		}
		fmt.Println("detail id: %s\n ", id.String())
	}
}
func InitDetailService[P any, R domain.Detail](
	ch *amqp.Channel, mongoRepo mongo_repo.IRepository[R],
	scraper func(url string) R) IService[P, R] {
	publish := rabbits.Publish[P](ch)
	return &serviceDetail[P, R]{publish: publish, scraper: scraper, mongoRepo: mongoRepo}
}
