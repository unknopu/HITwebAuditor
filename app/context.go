package app

import (
	"auditor/env"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

// Context app context
type Context struct {
	Environment *env.Environment
	RedisClient *redis.Client
	RedisLock   *redislock.Client
}

// NewContext new application context
func NewContext(e *env.Environment) *Context {
	// debug := e.Production == false
	return &Context{
		Environment: e,
	}
}
