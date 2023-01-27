package cron

import (
	cc "context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"

	"auditor/app"
	"auditor/env"
)

// JobInterface job interface
type JobInterface interface {
	Start()
}

// Job cron job
type Job struct {
	env    *env.Environment
	locker *locker
}

// locker implementation with Redis
type locker struct {
	cache *redis.Client
}

func (s *locker) Lock(key string) (success bool, err error) {
	ctx, cancel := cc.WithTimeout(cc.Background(), 2*time.Second)
	defer cancel()
	res, err := s.cache.SetNX(ctx, key, time.Now().String(), time.Second*5).Result()
	if err != nil {
		return false, err
	}
	return res, nil
}

func (s *locker) Unlock(key string) error {
	return nil
}

// NewCronJob new cron job
func NewCronJob(context *app.Context) JobInterface {
	job := &Job{
		env: context.Environment,
	}
	if context.Environment.Release {
		l := &locker{
			context.RedisClient,
		}
		job.locker = l
	}
	return job
}

// Start start
func (c *Job) Start() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	// defines a new scheduler that schedules and runs jobs
	schedule := gocron.NewScheduler(loc)

	// _, _ = schedule.Every(1).Hours().Do(c.TriggerCouponNotifications)

	if !c.env.Production {
		// _, _ = schedule.Every(5).Minutes().Do(c.TriggerPushOrderToReviewNotification)
	} else {
		// _, _ = schedule.Every(55).Minutes().Do(c.TriggerPushOrderToReviewNotification)
	}

	if c.env.Release {
		schedule.StartAsync()
	}
}
