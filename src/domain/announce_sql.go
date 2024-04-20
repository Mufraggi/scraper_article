package domain

import (
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AnnounceSql struct {
	Id             uuid.UUID
	AnnouncementId primitive.ObjectID
	DpeScore       string
	GesScore       string
	Characteristic []string
	Title          string
	Space          int
	Rooms          int
	City           string
	Price          float64
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
