package miss_configuration

import (
	"auditor/core/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "secMissCon"
)

// RepoInterface repo interface
type RepoInterface interface {
	Create(i interface{}) error
	Update(i interface{}) error
	Delete(i interface{}) error
	FindOneByID(id string, i interface{}) error
	FindOneByPrimitiveM(m primitive.M, i interface{}) error
	FindAllByPrimitiveM(m primitive.M, result interface{}, opts ...*options.FindOptions) error
}

// Repo otp repo
type Repo struct {
	mongodb.Repo
}

// NewRepo new service
func NewRepo() *Repo {
	return &Repo{
		Repo: mongodb.Repo{
			Collection: mongodb.
				DB().
				Collection(collectionName),
		},
	}
}
