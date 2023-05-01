package outdated_component

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
	Init(c *context.Context, f *OutdatedComponentForm) (interface{}, error)
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

func (s *Service) Init(c *context.Context, f *OutdatedComponentForm) (interface{}, error) {
	var reports []*entities.OutdatedComponentsReport

	for _, ref := range f.Refer {
		temp := strings.Split(ref, "/")
		if len(temp) < 1 {
			continue
		}

		if isPhpPWN(temp[1]) {
			v := entities.PhpPwnCVE()
			reports = append(reports, &entities.OutdatedComponentsReport{
				Location:       f.URL,
				Level:          entities.CRITICAL,
				Type:           entities.OutdatedComponents,
				Vaulnerability: entities.VULNERABILITY(v),
			})

		}
	}

	return buildPageInfomation(reports), nil
}
