package cryptograhpical_failure

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/core/utils"
	"auditor/entities"
	"sync"

	"github.com/jinzhu/copier"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *CFForm) (interface{}, error)
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

func (s *Service) Init(c *context.Context, f *CFForm) (interface{}, error) {
	option := f.URLOptions()

	bf := &utils.BasicRequestForm{}
	_ = copier.Copy(bf, option)
	response := utils.SendRequest(*bf, "")

	var reports []*entities.CryptoFailureReport
	if response.TLS == nil {
		reports = append(reports, &entities.CryptoFailureReport{
			Location:       f.URL,
			Level:          entities.MEDIUM,
			Type:           entities.Cryptography,
			Vaulnerability: entities.Certification,
		})
		reports = append(reports, &entities.CryptoFailureReport{
			Location:       f.URL,
			Level:          entities.HIGH,
			Type:           entities.Cryptography,
			Vaulnerability: entities.Transmittion,
		})
	}

	return buildPageInfomation(reports), nil
}
