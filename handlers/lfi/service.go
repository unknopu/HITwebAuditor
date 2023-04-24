package lfi

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	"regexp"
	"sync"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *LFIForm) (interface{}, error)
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

func (s *Service) Init(c *context.Context, f *LFIForm) (interface{}, error) {
	option := f.URLOptions()

	var p []string
	rule := `(?m)^[a-z_][a-z0-9_-]{0,30}[a-z0-9_$-]?:[^:]*:\d+:\d+:[^:]*:[^:]*:[^:]*$`
	for _, payload := range payloads {
		responseBody := injectPayload(*option, payload)

		isMatch, err := regexp.MatchString(rule, responseBody)
		if err != nil {
			return nil, err
		}

		if isMatch {
			p = append(p, payload)
		}
	}

	report := &entities.LFIReport{
		Location:       f.URL,
		Payload:        p,
		Level:          []string{"High"},
		Type:           entities.Broken,
		Vaulnerability: []entities.VULNERABILITY{entities.LocalFileIncusion},
	}

	return report, nil
}
