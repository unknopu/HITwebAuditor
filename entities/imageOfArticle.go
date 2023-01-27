package entities

import "auditor/core/mongodb"

type ImageOfArticle struct {
	mongodb.Model `bson:",inline"`
	FilePath      string `json:"imgfilepath,omitempty" bson:"img_file_path"`
}
