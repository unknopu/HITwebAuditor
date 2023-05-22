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
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *Form) (interface{}, error)
	GetLatest(c *context.Context, f *GetLatestForm) (interface{}, error)
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
	cache   *redis.Client
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
		cache:   c.RedisClient,
	}
}

func (s *Service) Init(c *context.Context, f *Form) (interface{}, error) {
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

	return nil, nil
}

func (s *Service) GetLatest(c *context.Context, f *GetLatestForm) (interface{}, error) {

	log.Println("init fetching... ")
	report, err := s.rp.FindLatest()
	if err != nil {
		return nil, err
	}

	sqli := s.fetchSQLI(c, *report.ID)
	secMissCon := s.fetchSecMissCon(c, *report.ID)
	outDateComponents := s.fetchOutdatedCpn(c, *report.ID)
	xss := s.fetchXSS(c, *report.ID)
	cryptoFailure := s.fetchCryptoFailure(c, *report.ID)
	lfi := s.fetchLFI(c, *report.ID)

	report.SQLi = sqli
	report.MConfig = secMissCon
	report.OutdatedComponent = outDateComponents
	report.XSS = xss
	report.CryptoFailure = cryptoFailure
	report.LFI = lfi
	return report, nil
}
