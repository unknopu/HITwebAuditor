package me

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/core/google"
	"auditor/entities"
	cc "context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	GetApplePublicKeys = "https://appleid.apple.com/auth/keys"
	AppleUrl           = "https://appleid.apple.com"
)

// ServiceInterface service interface
type ServiceInterface interface {
	TestDB(c *context.Context) (interface{}, error)
	TestRedis(c *context.Context) (interface{}, error)
}

// Service  repo
type Service struct {
	rp      RepoInterface
	c       *app.Context
	gs      google.SigninInterface
	context *app.Context
	redis   *redis.Client
}

// NewService new service
func NewService(c *app.Context) ServiceInterface {
	return &Service{
		c:       c,
		context: c,
		rp:      NewRepo(),
		redis:   c.RedisClient,
	}
}

func (s *Service) TestDB(c *context.Context) (interface{}, error) {
	uid := c.GetUserSession().UserID
	m := &entities.AccessToken{}
	err := s.rp.Create(m)
	if err != nil {
		return nil, err
	}
	return uid, nil
}

func (s *Service) TestRedis(c *context.Context) (interface{}, error) {
	ctx := cc.Background()
	val, err := s.redis.Get(ctx, "testXX").Result()
	if err != nil {
		val = "myData"

		time.Sleep(2 * time.Second)
		err := s.redis.Set(ctx, "testXX", val, 1*time.Minute).Err()
		if err != nil {
			return nil, err
		}
	}
	return fmt.Sprintf("[%v]", val), nil
}
