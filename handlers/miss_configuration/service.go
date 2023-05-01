package miss_configuration

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	"errors"
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
		return nil, errors.New("no leak")
	}

	var server, powerBy string
	if headerData.Server != "" {
		r := regexp.MustCompile(SERVER_RULE)
		server = r.FindString(headerData.Server)
	}
	if headerData.XPoweredBy != "" {
		r := regexp.MustCompile(POWERDBY_RULE)
		powerBy = r.FindString(headerData.XPoweredBy)
	}

	var reports []*entities.MissConfigurationReport
	if server != "" {
		if strings.ContainsAny(server, "nginxNginx") {
			reports = append(reports, &entities.MissConfigurationReport{
				Location:       f.URL,
				Payload:        []string{server},
				Level:          entities.LOW,
				Type:           entities.MisConfiguration,
				Vaulnerability: entities.NginxVersion,
			})
		}
		if strings.ContainsAny(server, "apacheApache") {
			reports = append(reports, &entities.MissConfigurationReport{
				Location:       f.URL,
				Payload:        []string{server},
				Level:          entities.LOW,
				Type:           entities.MisConfiguration,
				Vaulnerability: entities.ApacheVersion,
			})
		}
	}
	if powerBy != "" {
		reports = append(reports, &entities.MissConfigurationReport{
			Location:       f.URL,
			Payload:        []string{powerBy},
			Level:          entities.LOW,
			Type:           entities.MisConfiguration,
			Vaulnerability: entities.PHPVersion,
		})
	}

	return buildPageInfomation(reports), nil
}
