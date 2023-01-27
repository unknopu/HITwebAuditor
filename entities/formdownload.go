package entities

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductVideo video model
type FileForm struct {
	FileID      primitive.ObjectID `json:"id" bson:"_id"`
	FileName    string             `json:"file_name" bson:"file_name"`
	AccountName string             `json:"-" bson:"account_name"`
	Container   string             `json:"-" bson:"container"`
	FilePath    string             `json:"-" bson:"file_path"`
	CreateDt    time.Time          `json:"-" bson:"createDt"`
}

func (i FileForm) FileURL() string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", i.AccountName, i.Container, i.FilePath)
}

// MarshalJSON custom image json
func (i FileForm) MarshalJSON() ([]byte, error) {
	type Alias FileForm
	fileModel := &struct {
		*Alias
		FileURL string `json:"file_url"`
	}{
		Alias:   (*Alias)(&i),
		FileURL: i.FileURL(),
	}
	return json.Marshal(fileModel)
}
