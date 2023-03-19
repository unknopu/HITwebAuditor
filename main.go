package main

import (
	cc "context"
	"fmt"
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
	configPath := "configs"
	if os.Getenv("RELEASE") != "" {
		initConfigFile()
		configPath = "./"
	}

	envConfig, err := env.Read(configPath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	translator.InitTranslator()
	mongodbOptions := &mongodb.Options{
		URL:          envConfig.DatabaseURL,
		Port:         envConfig.DatabasePort,
		DatabaseName: envConfig.DatabaseName,
		Username:     envConfig.DatabaseUsername,
		Password:     envConfig.DatabasePassword,
		Root:         envConfig.DatabaseRoot,
		Debug:        !envConfig.Release,
	}
	err = mongodb.InitDatabase(mongodbOptions)
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

	server.New(router.NewWithOptions(options, context), envConfig.ServerPort).Start()
}

func initConfigFile() {
	f, err := os.Create("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	configs, err := os.OpenFile("config.yaml",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	d1 := (fmt.Sprintf("DATA_BASE_URL: %v\n", os.Getenv("DATA_BASE_URL")))
	_, _ = configs.WriteString(d1)

	dx := fmt.Sprintf("DATA_BASE_PORT: 27017\n")
	_, _ = configs.WriteString(dx)

	d2 := (fmt.Sprintf("DATA_BASE_NAME: %v\n", os.Getenv("DATA_BASE_NAME")))
	_, _ = configs.WriteString(d2)

	dr := (fmt.Sprintf("DATA_BASE_ROOT: true\n"))
	_, _ = configs.WriteString(dr)

	d3 := (fmt.Sprintf("DATA_BASE_USERNAME: %v\n", os.Getenv("DATA_BASE_USERNAME")))
	_, _ = configs.WriteString(d3)

	d4 := (fmt.Sprintf("DATA_BASE_PASSWORD: %v\n", os.Getenv("DATA_BASE_PASSWORD")))
	_, _ = configs.WriteString(d4)

	dport := fmt.Sprintf("SERVER_PORT: 8000\n")
	_, _ = configs.WriteString(dport)
}
