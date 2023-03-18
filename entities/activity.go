package entities

import (
	"auditor/core/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActivityType activity type
type ActivityType int

const (
	// ActivityTypeLogin activity type login
	ActivityTypeLogin ActivityType = iota + 1

	// ActivityTypeRegister activity type register
	ActivityTypeRegister

	// ActivityTypeRedeemPoint activity redeem point
	ActivityTypeRedeemPoint

	// ActivityTypeEditProfile activity edit profile
	ActivityTypeEditProfile
)

// Activity activity
type Activity struct {
	mongodb.Model `bson:",inline"`
	Name          string              `json:"name,omitempty" bson:"name"`
	Description   string              `json:"description,omitempty" bson:"description,omitempty"`
	Payload       interface{}         `json:"payload,omitempty" bson:"payload,omitempty"`
	UserID        *primitive.ObjectID `json:"user_id,omitempty" bson:"user_id"`
	Type          ActivityType        `json:"type" bson:"type"`
	ClientIP      string              `json:"client_ip,omitempty" bson:"client_ip,omitempty"`
	UserAgent     string              `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	ReferenceID   *primitive.ObjectID `json:"reference_id,omitempty" bson:"reference_id"`
}

type HealthCheck struct {
	mongodb.Model `bson:",inline"`
	Description   string `json:"description,omitempty" bson:"description,omitempty"`
	ClientIP      string `json:"client_ip,omitempty" bson:"client_ip,omitempty"`
	UserAgent     string `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
}
