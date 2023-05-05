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

	if isInjectable {
		if isErrMsgDetected {
			reports = append(reports, &entities.SQLiReport{
				Location:       f.URL,
				Level:          entities.CRITICAL,
				Type:           entities.Injection,
				Vaulnerability: entities.SQLIErr,
			})
		}
		if isContainBooleanBased(*option) {
			reports = append(reports, &entities.SQLiReport{
				Location:       f.URL,
				Level:          entities.CRITICAL,
				Type:           entities.Injection,
				Vaulnerability: entities.SQLIboolean,
			})
		}
	}

	if isContainUnionBased(*option) {
		reports = append(reports, &entities.SQLiReport{
			Location:       f.URL,
			Level:          entities.CRITICAL,
			Type:           entities.Injection,
			Vaulnerability: entities.SQLIUnion,
		})
	}

	return buildPageInfomation(reports), nil
}
