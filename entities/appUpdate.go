package entities

import "auditor/core/mongodb"

type EnumPlatform int

const (
	Android EnumPlatform = iota + 1
	IOS
)

type AppUpdate struct {
	mongodb.Model `bson:",inline"`
	Title         string `json:"title"       bson:"title,omitempty"`
	Message       string `json:"message"     bson:"message,omitempty"`
	BuildNumber   int    `json:"buildNumber" bson:"build_number,omitempty"`
	Force         bool   `json:"force"       bson:"force"`
	Platform      int    `json:"platform"    bson:"platform,omitempty"`
}
