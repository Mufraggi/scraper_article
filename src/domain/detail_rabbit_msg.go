package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type DetailRabbitMsg struct {
	Url     string             `json:"url"`
	MongoId primitive.ObjectID `json:"mongoId"`
}
