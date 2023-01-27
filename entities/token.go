package entities

import (
	"auditor/core/mongodb"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role role
type Role int

const (
	// RoleMember member
	RoleMember Role = iota
	// Member member
	Member
)

// AccessToken access token model
type AccessToken struct {
	mongodb.Model `bson:",inline"`
	UserID        primitive.ObjectID `json:"-" bson:"user_id"`
	Role          Role               `json:"role,omitempty" bson:"rele,omitempty"`
	AccessToken   string             `json:"access_token,omitempty" bson:"-"`
	RefreshToken  string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
}

// Claims jwt claims
type Claims struct {
	jwt.StandardClaims
	Type         int    `json:"type,omitempty"`
	Role         Role   `json:"role,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	SessionID    string `json:"session_id,omitempty"`
}

// ToUserSession convert claims to user session
func (c *Claims) ToUserSession() *UserContext {
	return &UserContext{
		UserID: c.Subject,
		Role:   c.Role,
	}
}

// UserContext user context
type UserContext struct {
	UserID         string
	Role           Role
	UniversityCode string
}

// GetObjectUserID get user id kide object id
func (u *UserContext) GetObjectUserID() primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(u.UserID)
	if err != nil {
		return primitive.NilObjectID
	}
	return oid
}
