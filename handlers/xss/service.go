package xss

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	"strings"
	"sync"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *XSSForm) (interface{}, error)
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

func (s *Service) Init(c *context.Context, f *XSSForm) (interface{}, error) {
	option := f.URLOptions()

	var p []string
	for _, payload := range payloads {
		bodyTag := fetchTagBody(*option, payload)
		if strings.ContainsAny(bodyTag, payload) {
			p = append(p, payload)
		}
		if len(p) > 5 {
			break
		}
	}

	report := &entities.XSSReport{
		Location:       f.URL,
		Payload:        p,
		Level:          []string{"High"},
		Type:           entities.Injection,
		Vaulnerability: []entities.VULNERABILITY{entities.CrossSiteScripting},
	}

	return report, nil
}
