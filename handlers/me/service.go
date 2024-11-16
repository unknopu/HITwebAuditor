package me

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/core/google"
	"auditor/entities"
	cc "context"
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
	m := &entities.HealthCheck{
		Description: "test mongodb",
		ClientIP:    c.RealIP(),
		UserAgent:   c.GetUserSession().UserAgent,
	}
	var msg string
	err := s.rp.Create(m)
	if err != nil {
		msg = err.Error()
	}
	rsp := &Healthcheck{Message: msg, IsHealthy: err == nil}
	return rsp, nil
}

func (s *Service) TestRedis(c *context.Context) (interface{}, error) {
	var msg string
	rsp := &Healthcheck{}

	ctx := cc.Background()
	if s.redis == nil {
		rsp.Message = "Redis servier not found"
		return rsp, nil
	}
	val, err := s.redis.Get(ctx, "testXX").Result()
	if err != nil {
		val = "myData"

		time.Sleep(2 * time.Second)
		err := s.redis.Set(ctx, "testXX", val, 1*time.Minute).Err()
		msg = err.Error()
	}

	rsp.Message = msg
	rsp.IsHealthy = true
	return rsp, nil
}
