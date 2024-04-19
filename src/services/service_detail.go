package services

import (
	"encoding/json"
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/Mufraggi/scraper_article/src/mongo_repo"
	"github.com/Mufraggi/scraper_article/src/rabbits"
	amqp "github.com/rabbitmq/amqp091-go"
)

type serviceDetail[P any] struct {
	publish   func(queueName string, msg P) error
	scraper   func(url string) domain.Detail
	mongoRepo mongo_repo.IRepository[domain.Detail]
}

func (s *serviceDetail[P]) GetProcessMsg() func(msg amqp.Delivery) {
	return func(msg amqp.Delivery) {
		var message domain.DetailRabbitMsg
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			fmt.Printf("error deserialization")
		}
		var res domain.Detail
		res = s.scraper(message.Url)
		id, err := s.mongoRepo.Create(res)
		if err != nil {
			fmt.Printf("error creating mongo repo")
		}
		fmt.Println("detail id: %s\n ", id.String())
	}
}
func InitDetailService[P any](
	ch *amqp.Channel, mongoRepo mongo_repo.IRepository[domain.Detail],
	scraper func(url string) domain.Detail) IService[P] {
	publish := rabbits.Publish[P](ch)
	return &serviceDetail[P]{publish: publish, scraper: scraper, mongoRepo: mongoRepo}
}
