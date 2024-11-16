package report

import (
	"auditor/core/mongodb"
	"auditor/entities"
	"context"
	"time"

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
	FindAllByPrimitiveM(m primitive.M, result interface{}, opts ...*options.FindOptions) error
	FindLatest() (*entities.Report, error)
	FindHistory() ([]*entities.Report, error)
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

func (r *Repo) FindHistory() ([]*entities.Report, error) {
	entity := []*entities.Report{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pipe := []primitive.M{
		primitive.M{
			"$sort": primitive.M{
				"updated_at": -1,
			},
		},
		primitive.M{
			"$skip": 0,
		},
		primitive.M{
			"$limit": 10,
		},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, nil
	}

	err = cursor.All(ctx, &entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
