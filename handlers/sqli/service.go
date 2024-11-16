package sqli

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *SqliForm) ([]*entities.SQLiReport, error)
	FetchReport(c *context.Context, id primitive.ObjectID) (*entities.Page, error)
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

func (s *Service) Init(c *context.Context, f *SqliForm) ([]*entities.SQLiReport, error) {
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

	// for index, _ := range reports {
	// 	reports[index].ReportNumber = f.ReportNumber
	// }

	// err := s.rp.Create(reports)
	// if err != nil {
	// 	return nil, err
	// }
	return reports, nil
}

func (s *Service) FetchReport(c *context.Context, id primitive.ObjectID) (*entities.Page, error) {
	reports := []*entities.SQLiReport{}
	err := s.rp.FindAllByPrimitiveM(primitive.M{"report_number": id}, &reports)
	if err != nil {
		return nil, err
	}

	return BuildPageInfomation(reports), nil

}
