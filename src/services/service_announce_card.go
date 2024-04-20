package services

import (
	"encoding/json"
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/Mufraggi/scraper_article/src/mongo_repo"
	"github.com/Mufraggi/scraper_article/src/rabbits"
	amqp "github.com/rabbitmq/amqp091-go"
)

type serviceAnnounceCard[P domain.DetailPublish] struct {
	publish   func(queueName string, msg domain.DetailPublish) error
	scraper   func(url string) (domain.ListAnnounceCard, error)
	mongoRepo mongo_repo.IRepository[domain.AnnounceCard]
}

func (s *serviceAnnounceCard[P]) GetProcessMsg() func(msg amqp.Delivery) {
	scraper := s.scraper
	return func(msg amqp.Delivery) {
		var message domain.DetailRabbitMsg
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			fmt.Printf("error deserialization")
		}

		res, err := scraper(message.Url)
		if err != nil {
			return
		}
		for _, v := range res {
			id, err := s.mongoRepo.Create(v)
			if err != nil {
				fmt.Printf("error creating mongo repo")
			}
			publishMsg := domain.DetailPublish{Id: id.String(), Url: "https://immobilier.lefigaro.fr" + v.URL}
			err = s.publish("detail", publishMsg)
			if err != nil {
				fmt.Printf("error during the message publish")
			}
			fmt.Println("detail id: %s\n ", id.String())
		}

	}
}
func InitAnnounceCardService[P domain.DetailPublish](
	ch *amqp.Channel,
	mongoRepo mongo_repo.IRepository[domain.AnnounceCard],
	scraper func(url string) (domain.ListAnnounceCard, error),
) IService[P] {
	publish := rabbits.Publish[domain.DetailPublish](ch)
	return &serviceAnnounceCard[P]{publish: publish, scraper: scraper, mongoRepo: mongoRepo}
}
