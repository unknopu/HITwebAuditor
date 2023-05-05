package report

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	cf "auditor/handlers/cryptograhpical_failure"
	lfi "auditor/handlers/lfi"
	mc "auditor/handlers/miss_configuration"
	odc "auditor/handlers/outdated_component"
	sqli "auditor/handlers/sqli"
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
	lfis    lfi.ServiceInterface
	sqlis   sqli.ServiceInterface
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
		lfis:    lfi.NewService(c),
		sqlis:   sqli.NewService(c),
	}
}

func (s *Service) Init(c *context.Context, f *Form) (interface{}, error) {
	// option := f.URLOptions()
	var missConfigVul, outdatedCpnVul, cryptoFailureVul, xssVul, lfiVul, sqliVul *entities.Page
	wg.Add(5)

	go func() {
		missConfigVul = s.doMissConfig(c, f)
		if missConfigVul != nil {
			outdatedCpnVul = s.doOutdatedCpn(c, missConfigVul.Entities.([]*entities.MissConfigurationReport))
		}
		wg.Done()
	}()
	go func() {
		cryptoFailureVul = s.doCryptoFailure(c, f)
		wg.Done()
	}()
	go func() {
		xssVul = s.doXSS(c, f)
		wg.Done()
	}()
	go func() {
		lfiVul = s.doLFI(c, f)
		wg.Done()
	}()
	go func() {
		sqliVul = s.doSQLI(c, f)
		wg.Done()
	}()
	wg.Wait()
	
	return &entities.Report{
		URL:               f.URL,
		MConfig:           missConfigVul,
		CryptoFailure:     cryptoFailureVul,
		OutdatedComponent: outdatedCpnVul,
		XSS:               xssVul,
		LFI:               lfiVul,
		SQLi:              sqliVul,
	}, nil
}
