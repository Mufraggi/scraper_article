package main

import (
	"context"
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/Mufraggi/scraper_article/src/listeners"
	"github.com/Mufraggi/scraper_article/src/mongo_repo"
	"github.com/Mufraggi/scraper_article/src/rabbits"
	"github.com/Mufraggi/scraper_article/src/scrapers/announce"
	"github.com/Mufraggi/scraper_article/src/services"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type TmpPublish struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
type DatabaseMongoConfig struct {
	MongoURI     string
	DatabaseName string
}

func InitDb(c DatabaseMongoConfig) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(c.MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	database := client.Database(c.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	return database, err
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("E RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erreur lors de la création du canal: %v", err)
	}
	defer ch.Close()

	db, _ := InitDb(DatabaseMongoConfig{
		MongoURI:     "mongodb://localhost:27017",
		DatabaseName: "testdb",
	})
	repoDetail := mongo_repo.InitDetailRepository(db)
	_ = rabbits.InitQueue(ch, "detail")
	consume, err := rabbits.Consume(ch, "detail")
	fmt.Println(err)
	getDetail := announce.GetAnnounceDetail()

	service := services.InitDetailService[TmpPublish, domain.Detail](ch, repoDetail, getDetail)
	process := service.GetProcessMsg()
	runDetail := listeners.InitListeners(consume, process).Run()
	go runDetail()
	//announce.GetAnnounceDetail("https://immobilier.lefigaro.fr/annonces/annonce-68148436.html")
	//url := "https://immobilier.lefigaro.fr/annonces/annonce-68271282.html"
	//url := "https://immobilier.lefigaro.fr/annonces/annonce-68271204.html"
	//	url := "https://immobilier.lefigaro.fr/annonces/annonce-68268986.html"

	fmt.Println("that run ")
	select {}
}
