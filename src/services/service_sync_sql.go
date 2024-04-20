package services

import (
	"encoding/json"
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/Mufraggi/scraper_article/src/mongo_repo"
	"github.com/Mufraggi/scraper_article/src/sql_repo"
	amqp "github.com/rabbitmq/amqp091-go"
)

type serviceSyncSql[P any] struct {
	mongoRepo        mongo_repo.IRepository[domain.Detail]
	announcementRepo sql_repo.IAdminRepository
}

func (s *serviceSyncSql[P]) GetProcessMsg() func(msg amqp.Delivery) {
	return func(msg amqp.Delivery) {
		var message domain.SyncAnnouncement
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			fmt.Printf("error deserialization")
		}
		announcement, err := s.mongoRepo.FindById(message.Id)
		if err != nil {
			fmt.Printf("error creating mongo repo")
		}
		id, err := s.announcementRepo.InsertAnnounce(announcement)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("announcement id: %s\n ", id.String())
	}
}
func InitSyncSqlService[P any](
	mongoRepo mongo_repo.IRepository[domain.Detail],
	announcementRepo sql_repo.IAdminRepository,
) IService[P] {
	return &serviceSyncSql[P]{
		announcementRepo: announcementRepo,
		mongoRepo:        mongoRepo,
	}
}
