package me

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/core/google"
)

type loginType int

const (
	facebookLogin loginType = iota + 1
	googleLogin
	lineLogin
	adminLogin
	appleLogin
)

const (
	GetApplePublicKeys = "https://appleid.apple.com/auth/keys"
	AppleUrl           = "https://appleid.apple.com"
)

// ServiceInterface service interface
type ServiceInterface interface {
	TEST(ctx *context.Context) (interface{}, error)
}

// Service  repo
type Service struct {
	c       *app.Context
	gs      google.SigninInterface
	context *app.Context
}

// NewService new service
func NewService(c *app.Context) ServiceInterface {
	return &Service{
		c:       c,
		context: c,
	}
}

func (s *Service) TEST(ctx *context.Context) (interface{}, error) {
	uid := ctx.GetUserSession().UserID
	return uid, nil
}
