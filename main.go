package main

import (
	cc "context"
	"os"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"

	"auditor/app"
	"auditor/core/mongodb"
	"auditor/core/server"
	"auditor/core/translator"
	_ "auditor/docs"
	"auditor/env"
	"auditor/logx"
	"auditor/middleware"
	"auditor/router"
)

func init() {

	//runtime.GOMAXPROCS(1)
}

// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs"
	}
	envConfig, _ := env.Read(configPath)

	translator.InitTranslator()
	mongodbOptions := validateMongoDB(envConfig)
	err := mongodb.InitDatabase(mongodbOptions)
	if err != nil {
		panic(err)
	}

	context := app.NewContext(envConfig)
	if os.Getenv("RedisOn") != "" {
		client := redis.NewClient(&redis.Options{
			Addr:     envConfig.RedisHost,
			Password: envConfig.RedisPassword,
		})
		ctx, cancel := cc.WithTimeout(cc.Background(), 2*time.Second)
		defer cancel()
		err = client.Ping(ctx).Err()
		if err != nil {
			panic(err)
		}
		context.RedisClient = client
		locker := redislock.New(client)
		context.RedisLock = locker
	}

	options := &router.Options{
		Environment: envConfig,
	}
	// if os.Getenv("RELEASE") != "" {
	// 	logx.Init("main", "trace")
	// 	options.LogLevel = log.INFO
	// 	options.LogMiddleware = middleware.Logger()
	// } else {
	logx.Init("main", "debug")
	options.LogLevel = log.DEBUG
	options.LogHeader = "\033[1;34m-->\033[0m ${time_rfc3339} ${level}"
	options.LogMiddleware = middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "\033[1;34m-->\033[0m method=${method} \033[1;32muri=${uri}\033[0m user_agent=${user_agent} " +
			"statu=${status} error=${error} latency_human=${latency_human}, \033[1;93mparameters=${parameters}\033[0m\n",
	})
	// }
	// schedule := cron.NewCronJob(context)
	// go schedule.Start()

	server.New(router.NewWithOptions(options, context), "8000").Start()
}

func validateMongoDB(envConfig *env.Environment) *mongodb.Options {

	if os.Getenv("RELEASE") == "" || envConfig == nil {
		return &mongodb.Options{
			URL:          envConfig.DatabaseURL,
			Port:         envConfig.DatabasePort,
			DatabaseName: envConfig.DatabaseName,
			Username:     envConfig.DatabaseUsername,
			Password:     envConfig.DatabasePassword,
			Root:         envConfig.DatabaseRoot,
			Debug:        !envConfig.Release,
		}
	}

	DATA_BASE_URL := os.Getenv("DATA_BASE_URL")
	DATA_BASE_NAME := os.Getenv("DATA_BASE_NAME")
	DATA_BASE_USERNAME := os.Getenv("DATA_BASE_USERNAME")
	DATA_BASE_PASSWORD := os.Getenv("DATA_BASE_PASSWORD")

	return &mongodb.Options{
		URL:          DATA_BASE_URL,
		Port:         27017,
		DatabaseName: DATA_BASE_NAME,
		Root:         true,
		Username:     DATA_BASE_USERNAME,
		Password:     DATA_BASE_PASSWORD,
		Debug:        true,
	}
}
