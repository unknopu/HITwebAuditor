package report

import (
	"auditor/core/mongodb"
	"auditor/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "report"
)

// RepoInterface repo interface
type RepoInterface interface {
	Create(i interface{}) error
	Update(i interface{}) error
	Delete(i interface{}) error
	FindLatest() (*entities.Report, error)
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

func (r *Repo) FindLatest() (*entities.Report, error) {
	entity := []*entities.Report{}
	o := options.
		Find().
		SetSort(bson.D{{"updated_at", -1}})
	err := r.FindAllByPrimitiveM(primitive.M{}, &entity, o)
	if err != nil {
		return nil, err
	}

	return entity[0], nil
}
