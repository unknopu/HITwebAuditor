package entities

import (
	"auditor/core/mongodb"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Image image model
type Image struct {
	mongodb.Model `bson:",inline"`
	UserID        primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	AccountName   string             `json:"-" bson:"account_name"`
	Container     string             `json:"-" bson:"container"`
	ThumbnailPath string             `json:"-" bson:"thumbnail_path"`
	LargePath     string             `json:"-" bson:"large_path"`
	MD5           string             `json:"md5,omitempty" bson:"md5"`
}

func (i Image) LargeURL() string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", i.AccountName, i.Container, i.LargePath)
}

func (i Image) ThumbnailURL() string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", i.AccountName, i.Container, i.ThumbnailPath)
}

// MarshalJSON custom image json
func (i Image) MarshalJSON() ([]byte, error) {
	type Alias Image
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

// ImageCompat Image path
type ImageCompat struct {
	AccountName string `json:"-" bson:"account_name"`
	Container   string `json:"-" bson:"container"`
	LargePath   string `json:"-" bson:"large_path"`
	ImagePath   string `json:"-" bson:"image_path"`
}

func (i ImageCompat) ImageURL() string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", i.AccountName, i.Container, i.LargePath)
}

// MarshalJSON custom image json
func (i ImageCompat) MarshalJSON() ([]byte, error) {
	type Alias ImageCompat
	imageModel := &struct {
		*Alias
		ImagePath string `json:"image_url,omitempty"`
	}{
		Alias:     (*Alias)(&i),
		ImagePath: i.ImageURL(),
	}
	return json.Marshal(imageModel)
}
