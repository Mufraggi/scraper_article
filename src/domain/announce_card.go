package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AnnounceCard struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	ProductType  string             `bson:"product_type"`
	Price        string             `bson:"price"`
	SquareFetter string             `bson:"square_fetter"`
	BeadRoom     string             `bson:"bead_room"`
	URL          string             `bson:"url"`
	Description  string             `bson:"description"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}
