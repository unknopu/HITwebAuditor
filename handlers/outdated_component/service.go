package outdated_component

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/core/utils"
	"auditor/entities"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *OutdatedComponentForm) ([]*entities.OutdatedComponentsReport, error)
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

func (s *Service) Init(c *context.Context, f *OutdatedComponentForm) ([]*entities.OutdatedComponentsReport, error) {
	var reports []*entities.OutdatedComponentsReport

	if len(f.Refer) == 0 {
		option := f.URLOptions()
		headerData := fetchHeaders(*option)
		if headerData.Server != "" {
			f.Refer = append(f.Refer, headerData.Server)
		}
		if headerData.XPoweredBy != "" {
			f.Refer = append(f.Refer, headerData.XPoweredBy)
		}
	}

	for _, ref := range f.Refer {
		temp := strings.Split(ref, "/")
		if len(temp) < 1 {
			continue
		}

		if isPhpPWN(temp[1]) && utils.IsExisting(temp[0], []string{"PHP", "php"}) {
			v := entities.PhpPwnCVE()
			reports = append(reports, &entities.OutdatedComponentsReport{
				Location:       f.URL,
				Level:          entities.CRITICAL,
				Type:           entities.OutdatedComponents,
				Vaulnerability: entities.VULNERABILITY(v),
			})

		}
	}

	for index, _ := range reports {
		reports[index].ReportNumber = f.ReportNumber
	}

	err := s.rp.Create(reports)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (s *Service) FetchReport(c *context.Context, id primitive.ObjectID) (*entities.Page, error) {
	reports := []*entities.OutdatedComponentsReport{}
	err := s.rp.FindAllByPrimitiveM(primitive.M{"report_number": id}, &reports)
	if err != nil {
		return nil, err
	}
	return BuildPageInfomation(reports), nil
}
