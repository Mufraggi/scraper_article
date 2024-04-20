package services

import (
	"encoding/json"
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/Mufraggi/scraper_article/src/mongo_repo"
	"github.com/Mufraggi/scraper_article/src/rabbits"
	amqp "github.com/rabbitmq/amqp091-go"
)

type serviceDetail[P domain.SyncAnnouncement] struct {
	publish   func(queueName string, msg domain.SyncAnnouncement) error
	scraper   func(url string) (*domain.Detail, error)
	mongoRepo mongo_repo.IRepository[domain.Detail]
}

func (s *serviceDetail[P]) GetProcessMsg() func(msg amqp.Delivery) {
	return func(msg amqp.Delivery) {
		var message domain.DetailRabbitMsg
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			fmt.Printf("error deserialization")
		}
		res, err := s.scraper(message.Url)
		if err != nil {
			return
		}
		id, err := s.mongoRepo.Create(*res)
		if err != nil {
			fmt.Printf("error creating mongo repo")
		}
		fmt.Println("detail id: %s\n ", id.String())
		messageSender := domain.SyncAnnouncement{
			Id: *id,
		}
		err = s.publish("sql_sync", messageSender)
		if err != nil {
			fmt.Printf("error during the message publish")
		}
		fmt.Println("detail id: %s\n ", id.String())
	}
}
func InitDetailService[P domain.SyncAnnouncement](
	ch *amqp.Channel, mongoRepo mongo_repo.IRepository[domain.Detail],
	scraper func(url string) (*domain.Detail, error)) IService[P] {
	publish := rabbits.Publish[domain.SyncAnnouncement](ch)
	return &serviceDetail[domain.SyncAnnouncement]{publish: publish, scraper: scraper, mongoRepo: mongoRepo}
}
