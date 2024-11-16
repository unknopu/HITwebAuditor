package response

import (
	"auditor/handlers/common"
	"encoding/json"
	"reflect"
)

// NewPage new page
func NewPage(entities interface{}) *Page {
	return &Page{
		PageSize: common.DefaultPageSize,
		Entities: entities,
	}
}

// NewPageWithSize new page with size
func NewPageWithSize(size int, entities interface{}) *Page {
	return &Page{
		PageSize: size,
		Entities: entities,
	}
}

// Page page model
type Page struct {
	PageSize         int         `json:"-"`
	CustomNextPageID interface{} `json:"-"`
	Entities         interface{} `json:"entities,omitempty"`
}

// SetPageSize set page size
func (page *Page) setPageSize(size int) *Page {
	page.PageSize = size
	return page
}

// SetFieldForNextPage set field
func (page *Page) setEntities(entities interface{}) *Page {
	page.Entities = entities
	return page
}

// PageInformation page information
type PageInformation struct {
	Size          int         `json:"size"`
	NumberOfItems int         `json:"number_of_items"`
	NextPageID    interface{} `json:"next_page_id,omitempty"`
}

// MarshalJSON custom page json
func (page Page) MarshalJSON() ([]byte, error) {
	type Alias Page
	pageModel := &struct {
		PageInfo *PageInformation `json:"page_information,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(&page),
	}

	pageSize := page.PageSize
	info := &PageInformation{}
	info.Size = pageSize
	v := reflect.ValueOf(page.Entities)
	numberOfItems := v.Len()
	if numberOfItems <= 0 {
		return json.Marshal(pageModel)
	}

	info.NumberOfItems = numberOfItems
	if numberOfItems == pageSize {
		if page.CustomNextPageID != nil {
			info.NextPageID = page.CustomNextPageID
		} else {
			itemValue := v.Index(pageSize - 1)
			idValue := reflect.Indirect(itemValue).FieldByName("ID")
			if idValue.IsValid() {
				info.NextPageID = idValue.Interface()
			}
		}
	}

	pageModel.PageInfo = info
	return json.Marshal(pageModel)
}
