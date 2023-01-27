package entities

import "auditor/core/mongodb"

type ImageOwnerBrand struct {
	mongodb.Model `bson:",inline"`
	FilePath      string `json:"imgfilepath" bson:"img_file_path"`
}
