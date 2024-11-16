package entities

import (
	"auditor/core/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Log log
type Log struct {
	mongodb.Model  `bson:",inline"`
	Uri            string              `json:"uri" bson:"uri"`
	ClientIP       string              `json:"client_ip,omitempty" bson:"client_ip,omitempty"`
	UserAgent      string              `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	Payload        interface{}         `json:"payload,omitempty" bson:"payload,omitempty"`
	HttpStatusCode int                 `json:"http_status_code,omitempty" bson:"http_status_code,omitempty"`
	UserID         *primitive.ObjectID `json:"user_id,omitempty" bson:"user_id"`
}
