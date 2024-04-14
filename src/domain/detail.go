package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Detail struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	AnnounceCardId primitive.ObjectID `bson:"announce_card_id"`
	DpeScore       string             `bson:"dpeScore"`
	GesScore       string             `bson:"gesScore"`
	Characteristic []string           `bson:"characteristic"`
	Title          string             `bson:"title"`
	Type           string             `bson:"type"`
	Space          string             `bson:"space"`
	Rooms          string             `bson:"rooms"`
	City           string             `bson:"city"`
	Price          string             `bson:"price"`
	Description    string             `bson:"description"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}
