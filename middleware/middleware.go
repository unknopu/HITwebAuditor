package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"auditor/logx"
)

type (
	// RateLimiterStore is the interface to be implemented by custom stores.
	RateLimiterStore interface {
		// Stores for the rate limiter have to implement the Allow method
		Allow(identifier string) (bool, error)
	}
)

type ManyRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(c echo.Context) bool
	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFunc func(c echo.Context)
	// RateLimiterConfig defines the configuration for the rate limiter
	RateLimiterConfig struct {
		Skipper    Skipper
		BeforeFunc BeforeFunc
		// IdentifierExtractor uses echo.Context to extract the identifier for a visitor
		IdentifierExtractor Extractor
		// Store defines a store for the rate limiter
		Store RateLimiterStore
		// ErrorHandler provides a handler to be called when IdentifierExtractor returns an error
		ErrorHandler func(context echo.Context, err error) error
		// DenyHandler provides a handler to be called when RateLimiter denies access
		DenyHandler func(context echo.Context, identifier string, err error) error
	}
	// Extractor is used to extract data from echo.Context
	Extractor func(context echo.Context) (string, error)
)

var (
	// DefaultSkipper default of skipper
	DefaultSkipper = func(c echo.Context) bool {
		blockPath := []string{"/builds", "/health", "/metrics"}
		for _, b := range blockPath {
			if c.Path() == b {
				return true
			}
		}
		return false
	}
	ExecTimeCache = make(map[string]time.Time)
)

var (
	// ErrRateLimitExceeded denotes an error raised when rate limit is exceeded
	ErrRateLimitExceeded = echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
	// ErrExtractorError denotes an error raised when extractor function is unsuccessful
	ErrExtractorError = echo.NewHTTPError(http.StatusForbidden, "error while extracting identifier")
)

// DefaultRateLimiterConfig defines default values for RateLimiterConfig
var DefaultRateLimiterConfig = RateLimiterConfig{
	Skipper: DefaultSkipper,
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		id := ctx.RealIP()
		return id, nil
	},
	ErrorHandler: func(context echo.Context, err error) error {
		return &echo.HTTPError{
			Code:     ErrExtractorError.Code,
			Message:  ErrExtractorError.Message,
			Internal: err,
		}
	},
	DenyHandler: func(context echo.Context, identifier string, err error) error {
		return &echo.HTTPError{
			Code:     ErrRateLimitExceeded.Code,
			Message:  ErrRateLimitExceeded.Message,
			Internal: err,
		}
	},
}

func (m *Middleware) RateLimiter(store RateLimiterStore) echo.MiddlewareFunc {
	config := DefaultRateLimiterConfig
	config.Store = store

	return m.RateLimiterWithConfig(config)
}

func (m *Middleware) RateLimiterWithConfig(config RateLimiterConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRateLimiterConfig.Skipper
	}
	if config.IdentifierExtractor == nil {
		config.IdentifierExtractor = DefaultRateLimiterConfig.IdentifierExtractor
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultRateLimiterConfig.ErrorHandler
	}
	if config.DenyHandler == nil {
		config.DenyHandler = DefaultRateLimiterConfig.DenyHandler
	}
	if config.Store == nil {
		logx.GetLog().Fatal("Store configuration must be provided")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}

			identifier, err := config.IdentifierExtractor(c)
			if err != nil {
				c.Error(config.ErrorHandler(c, err))
				return nil
			}
			if allow, _ := config.Store.Allow(identifier); !allow {
				return c.JSON(http.StatusTooManyRequests, ManyRequest{
					Code: 429, Message: "Too Many Requests",
				})

			}
			return next(c)
		}
	}
}

type (
	// RateLimiterMemoryStore is the built-in store implementation for RateLimiter
	RateLimiterMemoryStore struct {
		visitors map[string]*Visitor
		mutex    sync.Mutex
		rate     rate.Limit //for more info check out Limiter docs - https://pkg.go.dev/golang.org/x/time/rate#Limit.

		burst       int
		expiresIn   time.Duration
		lastCleanup time.Time
	}
	// Visitor signifies a unique user's limiter details
	Visitor struct {
		*rate.Limiter
		lastSeen time.Time
	}
)

func (m *Middleware) NewRateLimiterMemoryStore(rate rate.Limit) (store *RateLimiterMemoryStore) {
	return m.NewRateLimiterMemoryStoreWithConfig(RateLimiterMemoryStoreConfig{
		Rate: rate,
	})
}

func (m *Middleware) NewRateLimiterMemoryStoreWithConfig(config RateLimiterMemoryStoreConfig) (store *RateLimiterMemoryStore) {
	store = &RateLimiterMemoryStore{}

	store.rate = config.Rate
	store.burst = config.Burst
	store.expiresIn = config.ExpiresIn
	if config.ExpiresIn == 0 {
		store.expiresIn = DefaultRateLimiterMemoryStoreConfig.ExpiresIn
	}
	if config.Burst == 0 {
		store.burst = int(config.Rate)
	}
	store.visitors = make(map[string]*Visitor)
	store.lastCleanup = now()
	return
}

// RateLimiterMemoryStoreConfig represents configuration for RateLimiterMemoryStore
type RateLimiterMemoryStoreConfig struct {
	Rate      rate.Limit    // Rate of requests allowed to pass as req/s. For more info check out Limiter docs - https://pkg.go.dev/golang.org/x/time/rate#Limit.
	Burst     int           // Burst additionally allows a number of requests to pass when rate limit is reached
	ExpiresIn time.Duration // ExpiresIn is the duration after that a rate limiter is cleaned up
}

// DefaultRateLimiterMemoryStoreConfig provides default configuration values for RateLimiterMemoryStore
var DefaultRateLimiterMemoryStoreConfig = RateLimiterMemoryStoreConfig{
	ExpiresIn: 3 * time.Minute,
}

// Allow implements RateLimiterStore.Allow
func (store *RateLimiterMemoryStore) Allow(identifier string) (bool, error) {
	store.mutex.Lock()
	limiter, exists := store.visitors[identifier]
	if !exists {
		limiter = new(Visitor)
		limiter.Limiter = rate.NewLimiter(store.rate, store.burst)
		store.visitors[identifier] = limiter
	}
	limiter.lastSeen = now()
	if now().Sub(store.lastCleanup) > store.expiresIn {
		store.cleanupStaleVisitors()
	}
	store.mutex.Unlock()
	return limiter.AllowN(now(), 1), nil
}

/*
cleanupStaleVisitors helps manage the size of the visitors map by removing stale records
of users who haven't visited again after the configured expiry time has elapsed
*/
func (store *RateLimiterMemoryStore) cleanupStaleVisitors() {
	for id, visitor := range store.visitors {
		if now().Sub(visitor.lastSeen) > store.expiresIn {
			delete(store.visitors, id)
		}
	}
	store.lastCleanup = now()
}

/*
actual time method which is mocked in test file
*/
var now = time.Now

func (m *Middleware) GetExecTimeCache() map[string]time.Time {
	return ExecTimeCache
}

func (m *Middleware) GetStartTimeFromCache(c echo.Context, key string) time.Time {
	startTime := time.Now()
	if v, ok := ExecTimeCache[key]; ok {
		startTime = v
	}

	return startTime
}

const (
	requestInfoMsg       = "request information"
	responseInfoMsg      = "response information"
	logKeywordDontChange = "api_summary"
)

// Skipper skip middleware

type Middleware struct {
	Service  string
	Skipper  Skipper
	PubTopic string
}

func New(service string, args ...interface{}) *Middleware {
	m := &Middleware{
		Service: service,
		Skipper: DefaultSkipper,
	}

	return m
}

// Logger log request information
func (m *Middleware) Logger() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if m.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			logx.GetLog().WithFields(
				logrus.Fields{
					"@class":     "middelware",
					"@method":    req.Method,
					"@path_uri":  req.RequestURI,
					"@remote_ip": c.RealIP(),
					"@status":    res.Status,
					"@duration":  fmt.Sprint(stop.Sub(start).Milliseconds()) + "ms",
				},
			).Info()
			return
		}
	}
}

// Logging when have request to this service
func (m *Middleware) LogRequestInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logrus.SetFormatter(&logrus.JSONFormatter{})
			if m.Skipper(c) {
				return next(c)
			}

			reqBody := []byte{}
			if c.Request().Body != nil {
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			var requestMap map[string]interface{}
			_ = json.Unmarshal(reqBody, &requestMap)

			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

			jsonRequest, _ := json.Marshal(requestMap)
			logx.GetLog().WithFields(
				logrus.Fields{
					"@class":   "middelware",
					"@header":  c.Request().Header,
					"@request": string(jsonRequest),
				},
			).Trace()

			return next(c)
		}
	}
}

// Logging when service has response to requester.
func (m *Middleware) LogResponseInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			defer delete(ExecTimeCache, c.Request().Header.Get("spanId"))
			logrus.SetFormatter(&logrus.JSONFormatter{})
			if m.Skipper(c) {
				return next(c)
			}

			reqBody := []byte{}
			if c.Request().Body != nil {
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

			resBody := new(bytes.Buffer)

			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			if err := next(c); err != nil {
				c.Error(err)
			}

			var responseMap map[string]interface{}
			_ = json.Unmarshal(resBody.Bytes(), &responseMap)

			jsonResponse, _ := json.Marshal(responseMap)

			if len(string(jsonResponse)) <= 25600 {
				logx.GetLog().WithFields(logrus.Fields{"@class": "middelware", "@response": string(jsonResponse)}).Trace()
			} else {
				logx.GetLog().WithFields(logrus.Fields{"@class": "middelware", "@response": ""}).Trace()
			}

			return nil
		}
	}
}
func (m *Middleware) Build(buildstamp, githash string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/builds" {
				return c.JSON(http.StatusOK, map[string]string{
					"buildstamp": buildstamp,
					"githash":    githash,
				})
			}
			return next(c)
		}
	}
}
