package main

import (
	"context"
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/Mufraggi/scraper_article/src/listeners"
	"github.com/Mufraggi/scraper_article/src/mongo_repo"
	"github.com/Mufraggi/scraper_article/src/rabbits"
	"github.com/Mufraggi/scraper_article/src/scrapers/announce"
	"github.com/Mufraggi/scraper_article/src/scrapers/list_announce"
	"github.com/Mufraggi/scraper_article/src/services"
	"github.com/Mufraggi/scraper_article/src/sql_repo"
	"github.com/jackc/pgx/v4/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

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
	pgPool := getPool()
	repoDetail := mongo_repo.InitDetailRepository(db)
	repoListAnnounce := mongo_repo.InitAnnounceCardRepository(db)
	announcemntSqlRepo := sql_repo.InitAdminRepository(pgPool)

	_ = rabbits.InitQueue(ch, "list_announce")
	_ = rabbits.InitQueue(ch, "detail")
	_ = rabbits.InitQueue(ch, "sql_sync")

	consume, _ := rabbits.Consume(ch, "detail")
	consumerListAnnounce, _ := rabbits.Consume(ch, "list_announce")
	syncConsumer, _ := rabbits.Consume(ch, "sql_sync")

	getDetail := announce.GetAnnounceDetail()
	listAnnounce := list_announce.GetListAnnounce()

	serviceDetail := services.InitDetailService[domain.SyncAnnouncement](ch, repoDetail, getDetail)
	process := serviceDetail.GetProcessMsg()
	runDetail := listeners.InitListeners(consume, process).Run()
	go runDetail()

	serviceListAnnounce := services.InitAnnounceCardService[domain.DetailPublish](
		ch, repoListAnnounce, listAnnounce)
	processListAnnounce := serviceListAnnounce.GetProcessMsg()
	runList := listeners.InitListeners(consumerListAnnounce, processListAnnounce).Run()
	go runList()

	syncService := services.InitSyncSqlService[any](repoDetail, announcemntSqlRepo)
	processSync := syncService.GetProcessMsg()
	runSync := listeners.InitListeners(syncConsumer, processSync).Run()
	go runSync()

	//announce.GetAnnounceDetail("https://immobilier.lefigaro.fr/annonces/annonce-68148436.html")
	//url := "https://immobilier.lefigaro.fr/annonces/annonce-68271282.html"
	//url := "https://immobilier.lefigaro.fr/annonces/annonce-68271204.html"
	//	url := "https://immobilier.lefigaro.fr/annonces/annonce-68268986.html"

	fmt.Println("that run ")
	select {}
}

func getPool() *pgxpool.Pool {
	host := "localhost" // ou "127.0.0.1"
	port := 5432
	username := "user"
	password := "password"
	dbName := "gdg"
	sslmode := "disable"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbName, sslmode)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil
	}
	pool, _ := pgxpool.ConnectConfig(context.Background(), config)
	return pool
}
