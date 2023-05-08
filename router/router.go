package router

import (
	"auditor/handlers/errors"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"golang.org/x/time/rate"

	"auditor/app"
	"auditor/core/validator"
	"auditor/env"
	cf "auditor/handlers/cryptograhpical_failure"
	odc "auditor/handlers/outdated_component"
	sa "auditor/handlers/static_analysis"

	"auditor/handlers/lfi"
	"auditor/handlers/me"
	mc "auditor/handlers/miss_configuration"
	"auditor/handlers/report"
	"auditor/handlers/sql"
	"auditor/handlers/sqli"
	"auditor/handlers/xss"

	middlewareLog "auditor/middleware"
	myMiddleware "auditor/middleware"
)

var (
	buildstamp string
	githash    string
)

// Options option for new router
type Options struct {
	LogLevel      log.Lvl
	LogHeader     string
	LogMiddleware echo.MiddlewareFunc
	Environment   *env.Environment
}

func initEcho(m *middlewareLog.Middleware) *echo.Echo {
	e := echo.New()

	e.Use(m.Build(buildstamp, githash))
	if true {
		config := middlewareLog.RateLimiterConfig{
			Store: m.NewRateLimiterMemoryStore(rate.Limit(20)),
		}
		e.Use(m.RateLimiterWithConfig(config))
	}
	e.Use(m.LogRequestInfo())
	e.Use(m.Logger())
	return e
}

// New new router
func New() *echo.Echo {
	return NewWithOptions(nil, nil)
}

// NewWithOptions new router with options
func NewWithOptions(options *Options, context *app.Context) *echo.Echo {
	bundle := i18n.NewBundle(language.BritishEnglish)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("active.th.toml")
	bundle.MustLoadMessageFile("active.en.toml")
	m := middlewareLog.New("boot-api")
	router := initEcho(m)
	router.Validator = validator.New()
	router.HTTPErrorHandler = errors.HTTPErrorHandler

	router.Logger.SetPrefix("BOOTS")

	api := router.Group("/api/:version")
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"mobile",
			"session",
		},
	}))
	api.Use(middleware.Recover())
	api.Use(myMiddleware.CustomContext(bundle))
	if options != nil {
		router.Logger.SetLevel(options.LogLevel)
		if options.LogHeader != "" {
			router.Logger.SetHeader(options.LogHeader)
		}
		api.Use(options.LogMiddleware)
	}
	api.Use(
		middleware.Secure(),
		middleware.Gzip(),
		myMiddleware.ActivityLog(),
	)

	// API health checker
	api.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy\n")
	})

	ReportHandler := report.NewHandler(context)
	meHandler := me.NewHandler(context)
	SQLIHandler := sql.NewHandler(context)
	// spiderHandler := spider.NewHandler(context)
	SqliHandler := sqli.NewHandler(context)
	LFIHandler := lfi.NewHandler(context)
	XSSHandler := xss.NewHandler(context)
	MissConfigHandler := mc.NewHandler(context)
	CFHandler := cf.NewHandler(context)
	OutdatedHandler := odc.NewHandler(context)
	StaticAnalysisHandler := sa.NewHandler(context)

	meGroup := api.Group("/me")
	{
		meGroup.GET("/mongodb", meHandler.TestDB)
		meGroup.GET("/redis", meHandler.TestRedis)
		meGroup.GET("/healthcheck", meHandler.HealthCheck)
	}

	ReportGroup := api.Group("/report")
	{
		ReportGroup.POST("", ReportHandler.Init)
	}

	StaticAnalysisGroup := api.Group("/sa")
	{
		StaticAnalysisGroup.POST("", StaticAnalysisHandler.Init)
	}

	CFGroup := api.Group("/cf")
	{
		CFGroup.POST("", CFHandler.Init)
	}

	SQLiGroup := api.Group("/sqli")
	{
		SQLiGroup.GET("/test", SQLIHandler.TestIntruder)
		SQLiGroup.POST("", SQLIHandler.Init)
		SQLiGroup.POST("/error", SQLIHandler.ErrorBased)
		// SQLiGroup.POST("/union", SQLIHandler.UnionBased)

		// ==================================== v2
		SQLiGroup.POST("/start", SqliHandler.Init)
	}

	LFIGroup := api.Group("/lfi")
	{
		LFIGroup.POST("/start", LFIHandler.Init)
	}

	XSSGroup := api.Group("/xss")
	{
		XSSGroup.POST("/start", XSSHandler.Init)
	}

	MissConfigGroup := api.Group("/mc")
	{
		MissConfigGroup.POST("/start", MissConfigHandler.Init)
	}

	OutdatedComponentGroup := api.Group("/odc")
	{
		OutdatedComponentGroup.POST("", OutdatedHandler.Init)
	}

	return router
}
