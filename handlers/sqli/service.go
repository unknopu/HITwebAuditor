package sqli

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	"sync"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *SqliForm) (interface{}, error)
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

func (s *Service) Init(c *context.Context, f *SqliForm) (interface{}, error) {
	var reports []*entities.SQLiReport
	option := f.URLOptions()
	isInjectable, body := isParamInjectable(*option)
	isErrMsgDetected := detectErrMsg(body)

	wg.Add(3)
	if isInjectable {
		go func() {
			if isErrMsgDetected {
				reports = append(reports, &entities.SQLiReport{
					Location:       f.URL,
					Level:          entities.CRITICAL,
					Type:           entities.Injection,
					Vaulnerability: entities.SQLIErr,
				})
			}
			wg.Done()
		}()

		go func() {
			if isContainBooleanBased(*option) {
				reports = append(reports, &entities.SQLiReport{
					Location:       f.URL,
					Level:          entities.CRITICAL,
					Type:           entities.Injection,
					Vaulnerability: entities.SQLIboolean,
				})
			}
			wg.Done()
		}()

	}

	go func() {
		if isContainUnionBased(*option) {
			reports = append(reports, &entities.SQLiReport{
				Location:       f.URL,
				Level:          entities.CRITICAL,
				Type:           entities.Injection,
				Vaulnerability: entities.SQLIUnion,
			})
		}
		wg.Done()
	}()
	wg.Wait()

	return buildPageInfomation(reports), nil
}
