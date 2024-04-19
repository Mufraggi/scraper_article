package mongo_repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type AnnounceCardRepository struct {
	collection *mongo.Collection
}

func InitAnnounceCardRepository(db *mongo.Database) IRepository[domain.AnnounceCard] {
	c := db.Collection("AnnounceCard")
	return &AnnounceCardRepository{collection: c}
}

func (c *AnnounceCardRepository) Create(detail domain.AnnounceCard) (*primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	insertResult, err := c.collection.InsertOne(ctx, detail)

	if err != nil {
		log.Println("Error inserting into database:", err)
		return nil, err
	}
	id := insertResult.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (c *AnnounceCardRepository) FindById(id primitive.ObjectID) (*domain.AnnounceCard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var detail domain.AnnounceCard
	filter := bson.M{"_id": id}

	err := c.collection.FindOne(ctx, filter).Decode(&detail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("Company not found")
		} else {
			fmt.Println(err)
		}
		return nil, err
	}
	return &detail, err
}
