package report

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	cf "auditor/handlers/cryptograhpical_failure"
	mc "auditor/handlers/miss_configuration"
	odc "auditor/handlers/outdated_component"
	xss "auditor/handlers/xss"
	"sync"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *Form) (interface{}, error)
}

// Service  repo
type Service struct {
	rp      RepoInterface
	context *app.Context
	mcs     mc.ServiceInterface
	cfs     cf.ServiceInterface
	odcs    odc.ServiceInterface
	xsss    xss.ServiceInterface
}

// NewService new service
func NewService(c *app.Context) ServiceInterface {
	return &Service{
		context: c,
		rp:      NewRepo(),
		mcs:     mc.NewService(c),
		cfs:     cf.NewService(c),
		odcs:    odc.NewService(c),
		xsss:    xss.NewService(c),
	}
}

func (s *Service) Init(c *context.Context, f *Form) (interface{}, error) {
	// option := f.URLOptions()

	missConfig := s.doMissConfig(c, f)
	outdatedCpn := s.doOutdatedCpn(c, missConfig.Entities.([]*entities.MissConfigurationReport))
	cryptoFailure := s.doCryptoFailure(c, f)
	xssVul := s.doXSS(c, f)

	return &entities.Report{
		URL:               f.URL,
		MConfig:           missConfig,
		CryptoFailure:     cryptoFailure,
		OutdatedComponent: outdatedCpn,
		XSS:               xssVul,
	}, nil
}
