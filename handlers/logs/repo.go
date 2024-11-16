package logs

import (
	"auditor/core/mongodb"
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "logs"
)

// RepoInterface repo interface
type RepoInterface interface {
	Create(i interface{}) error
	Update(i interface{}) error
	Delete(i interface{}) error
	FindOneByID(id string, i interface{}) error
	FindAllWithAdminParameters(f *GetAllWithAdminForm, i interface{}) (int64, error)
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

// FindAllWithAdminParameters find all with parameters
func (r *Repo) FindAllWithAdminParameters(f *GetAllWithAdminForm, i interface{}) (int64, error) {
	uID, err := primitive.ObjectIDFromHex(f.ID)
	if err != nil {
		return 0, err
	}
	m := primitive.M{
		"user_id": uID,
	}
	size := f.PageSize()
	o := options.
		Find().
		SetSort(bson.M{
			"_id": -1,
		}).
		SetLimit(int64(size)).
		SetSkip(int64(size * (f.GetPage() - 1)))
	if err := r.FindAllByPrimitiveM(m, i, o); err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c, err := r.Collection.CountDocuments(ctx, m)
	if err != nil {
		return 0, mongodb.WrapError(err)
	}
	return int64(math.Ceil(float64(c) / float64(size))), nil
}
