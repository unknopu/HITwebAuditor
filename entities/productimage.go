package entities

import (
	"auditor/core/mongodb"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductImage image model
type ProductImage struct {
	mongodb.Model `bson:",inline"`
	ProductID     primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	AccountName   string             `json:"-" bson:"account_name"`
	Container     string             `json:"-" bson:"container"`
	ThumbnailPath string             `json:"-" bson:"thumbnail_path"`
	LargePath     string             `json:"-" bson:"large_path"`
}

func (i ProductImage) LargeURL() string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", i.AccountName, i.Container, i.LargePath)
}

func (i ProductImage) ThumbnailURL() string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", i.AccountName, i.Container, i.ThumbnailPath)
}

// MarshalJSON custom image json
func (i ProductImage) MarshalJSON() ([]byte, error) {
	type Alias ProductImage
	imageModel := &struct {
		*Alias
		ThumbnailURL string `json:"thumbnail_url,omitempty"`
		LargeURL     string `json:"large_url,omitempty"`
	}{
		Alias:        (*Alias)(&i),
		ThumbnailURL: i.ThumbnailURL(),
		LargeURL:     i.LargeURL(),
	}
	return json.Marshal(imageModel)
}
