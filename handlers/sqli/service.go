package sqli

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	"errors"
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
	option := f.URLOptions()
	IsPwn := tryInjection(*option)

	if IsPwn {
		return &entities.SQLiReport{
			Location:       f.URL,
			Level:          []string{"Critical"},
			Type:           entities.Injection,
			Vaulnerability: []entities.VULNERABILITY{entities.SQLIboolean},
		}, nil
	}

	return nil, errors.New("not found")
}
