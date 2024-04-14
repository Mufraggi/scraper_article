package mongo_repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository[T any] interface {
	Create(detail T) (*primitive.ObjectID, error)
	FindById(id primitive.ObjectID) (*T, error)
}
