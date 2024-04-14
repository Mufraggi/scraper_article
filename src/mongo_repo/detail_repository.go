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

type DetailRepository struct {
	collection *mongo.Collection
}

func InitDetailRepository(db *mongo.Database) IRepository[domain.Detail] {
	c := db.Collection("detail")
	return &DetailRepository{collection: c}
}

func (c *DetailRepository) Create(detail domain.Detail) (*primitive.ObjectID, error) {
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

func (c *DetailRepository) FindById(id primitive.ObjectID) (*domain.Detail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var detail domain.Detail
	filter := bson.M{"_id": id}

	err := c.collection.FindOne(ctx, filter).Decode(&detail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("Company not found")
		} else {
			log.Fatal(err)
		}
		return nil, err
	}
	return &detail, err
}
