package miss_configuration

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	"log"
	"regexp"
	"strings"
	"sync"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *MCForm) (interface{}, error)
}

// Service  repo
type Service struct {
	rp      RepoInterface
	context *app.Context
}

// NewService new service
func NewService(c *app.Context) ServiceInterface {
	return &Service{
		context: c,
		rp:      NewRepo(),
	}
}

func (s *Service) Init(c *context.Context, f *MCForm) (interface{}, error) {
	option := f.URLOptions()
	headerData := fetchHeaders(*option)

	if !anyVersionLeak(headerData) {
		return nil, nil
	}

	var server, powerBy string
	var phpPWN, serverPWN bool
	if headerData.Server != "" {
		r := regexp.MustCompile(SERVER_RULE)
		server = r.FindString(headerData.Server)
		serverPWN = false
		// serverPWN = checkPWNVersion(strings.Split(server, "/")[1], 1, 20)
	}
	if headerData.XPoweredBy != "" {
		r := regexp.MustCompile(POWERDBY_RULE)
		powerBy = r.FindString(headerData.XPoweredBy)
		phpPWN = checkPWNVersion("7.1.0", 5, 3, 7, 0)
	}

	sv := strings.Split(server, "/")
	pw := "7.1.0"
	log.Println()
	log.Println(sv, serverPWN)
	log.Println(pw, phpPWN)
	log.Println()

	report := &entities.MissConfigurationReport{
		Location:       f.URL,
		Payload:        []string{server, powerBy},
		Level:          []string{"LOW"},
		Type:           entities.MisConfiguration,
		Vaulnerability: []entities.VULNERABILITY{},
	}

	return report, nil
}
