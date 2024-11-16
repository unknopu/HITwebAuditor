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
	"errors"
	"log"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *Form) (interface{}, error)
	GetLatest(c *context.Context, f *GetLatestForm) (interface{}, error)
	GetHistory(c *context.Context, f *GetLatestForm) (interface{}, error)

	mapReports(c *context.Context, reports []*entities.Report)
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
	cache   *cache.Cache
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
		cache:   cache.New(1*time.Minute, 3*time.Minute),
	}
}

func (s *Service) Init(c *context.Context, f *Form) (interface{}, error) {
	_, hit := s.cache.Get(f.URL)
	if hit {
		return nil, errors.New("url target is processing.")
	}

	report := &entities.Report{URL: f.URL}
	err := s.rp.Create(report)
	if err != nil {
		return nil, err
	}
	f.ReportNumber = *report.ID

	wg.Add(5)

	go func() {
		missConfigVul := s.doMissConfig(c, f)
		if missConfigVul != nil {
			_ = s.doOutdatedCpn(c, missConfigVul, *report.ID)
		}
		wg.Done()
	}()
	go func() {
		_ = s.doCryptoFailure(c, f)
		wg.Done()
	}()
	go func() {
		_ = s.doXSS(c, f)
		wg.Done()
	}()
	go func() {
		_ = s.doLFI(c, f)
		wg.Done()
	}()
	go func() {
		_ = s.doSQLI(c, f)
		wg.Done()
	}()
	wg.Wait()

	s.cache.Set(f.URL, nil, 5*time.Second)
	return nil, nil
}

func (s *Service) GetLatest(c *context.Context, f *GetLatestForm) (interface{}, error) {
	log.Println("init fetching... ")
	report, err := s.rp.FindLatest()
	if err != nil {
		return nil, err
	}

	cache, hit := s.cache.Get(report.ID.Hex())
	if hit {
		return cache.(*entities.Report), nil
	}

	s.mapReports(c, []*entities.Report{report})
	s.cache.Set(report.ID.Hex(), report, 5*time.Second)

	return report, nil
}
func (s *Service) GetHistory(c *context.Context, f *GetLatestForm) (interface{}, error) {
	reports, err := s.rp.FindHistory()
	if err != nil {
		return nil, err
	}

	cache, hit := s.cache.Get(c.Request().URL.String())
	if hit {
		return cache.([]*entities.Report), nil
	}

	s.mapReports(c, reports)
	s.cache.Set(c.Request().URL.String(), reports, 35*time.Second)

	return reports, nil
}

func (s *Service) mapReports(c *context.Context, reports []*entities.Report) {
	for index, report := range reports {
		sqli := s.fetchSQLI(c, *report.ID)
		secMissCon := s.fetchSecMissCon(c, *report.ID)
		outDateComponents := s.fetchOutdatedCpn(c, *report.ID)
		xss := s.fetchXSS(c, *report.ID)
		cryptoFailure := s.fetchCryptoFailure(c, *report.ID)
		lfi := s.fetchLFI(c, *report.ID)

		reports[index].SQLi = sqli
		reports[index].MConfig = secMissCon
		reports[index].OutdatedComponent = outDateComponents
		reports[index].XSS = xss
		reports[index].CryptoFailure = cryptoFailure
		reports[index].LFI = lfi
	}
}
