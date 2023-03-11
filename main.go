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
	"auditor/core/utils"
	"auditor/docs"
	_ "auditor/docs"
	"auditor/entities"
	"auditor/env"
	"auditor/logx"
	"auditor/middleware"
	"auditor/response"
	"auditor/router"
)

func init() {

	//runtime.GOMAXPROCS(1)
}

// @title Boot Mobile API DOCS
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

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
	envConfig, err := env.Read(configPath)
	if err != nil {
		panic(err)
	}
	r, err := response.ReadReturnResult(configPath, "return_results")
	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.Host = envConfig.SwaggerHost
	docs.SwaggerInfo.BasePath = envConfig.SwaggerBasePath

	translator.InitTranslator()
	utils.InitSecureCookie(envConfig.CookieHashKey, envConfig.CookieBlockKey)

	entities.BaseURL = envConfig.WebURL
	databaseURL := envConfig.DatabaseURL
	if os.Getenv("IS_CLOUD") == "true" {
		databaseURL = "mongodb"
	}
	err = mongodb.InitDatabase(&mongodb.Options{
		URL:          databaseURL,
		Port:         envConfig.DatabasePort,
		DatabaseName: envConfig.DatabaseName,
		Username:     envConfig.DatabaseUsername,
		Password:     envConfig.DatabasePassword,
		Root:         envConfig.DatabaseRoot,
		Debug:        !envConfig.Release,
	})
	if err != nil {
		panic(err)
	}

	context := app.NewContext(envConfig, r)
	if envConfig.RedisOn {
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
		Results:     r,
	}
	if envConfig.Release {
		logx.Init("main", "trace")
		options.LogLevel = log.INFO
		options.LogMiddleware = middleware.Logger()
	} else {
		logx.Init("main", "debug")
		options.LogLevel = log.DEBUG
		options.LogHeader = "\033[1;34m-->\033[0m ${time_rfc3339} ${level}"
		options.LogMiddleware = middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "\033[1;34m-->\033[0m method=${method} \033[1;32muri=${uri}\033[0m user_agent=${user_agent} " +
				"statu=${status} error=${error} latency_human=${latency_human}, \033[1;93mparameters=${parameters}\033[0m\n",
		})
	}
	// schedule := cron.NewCronJob(context)
	// go schedule.Start()

	server.New(router.NewWithOptions(options, context), envConfig.ServerPort).Start()
}
