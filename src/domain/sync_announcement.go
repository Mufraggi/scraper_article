package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SyncAnnouncement struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}
