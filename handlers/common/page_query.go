package common

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DefaultPageSize default page size
var DefaultPageSize = 20

// PageQuery page form
type PageQuery struct {
	Query      string `query:"q"`
	NextPageID string `query:"next_page_id"`
	Size       int    `query:"size"`
}

// AddNextPage add next page
func (form *PageQuery) AddNextPage(m primitive.M) {
	id, err := primitive.ObjectIDFromHex(form.NextPageID)
	if err != nil {
		return
	}
	m["_id"] = primitive.M{
		"$gt": id,
	}
}

// AddReverse add reverse next page
func (form *PageQuery) AddReverse(m primitive.M) {
	id, err := primitive.ObjectIDFromHex(form.NextPageID)
	if err != nil {
		return
	}
	m["_id"] = primitive.M{
		"$lt": id,
	}
}

// PageSize page size
func (form *PageQuery) PageSize() int {
	if form.Size > 0 {
		return form.Size
	}
	return DefaultPageSize
}
