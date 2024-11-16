package common

import "go.mongodb.org/mongo-driver/bson/primitive"

// AdminPageQuery page form
type AdminPageQuery struct {
	Query string `json:"q,omitempty" query:"q"`
	Page  int    `query:"page"`
	Size  int    `query:"size"`
	NextPageID string `query:"next_page_id"`
}

// GetPage get page
func (form *AdminPageQuery) GetPage() int {
	if form.Page > 0 {
		return form.Page
	}
	return 1
}

// PageSize page size
func (form *AdminPageQuery) PageSize() int {
	if form.Size > 0 {
		return form.Size
	}
	return DefaultPageSize
}

// AddNextPage add next page
func (form *AdminPageQuery) AddNextPage(m primitive.M) {
	id, err := primitive.ObjectIDFromHex(form.NextPageID)
	if err != nil {
		return
	}
	m["_id"] = primitive.M{
		"$gt": id,
	}
}
