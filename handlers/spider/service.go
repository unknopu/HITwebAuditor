package spider

import (
	"auditor/app"
	"auditor/core/context"
	"fmt"
	"sync"
)

const (
	base = "/sql"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Start(c *context.Context, f *BaseForm) (interface{}, error)
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

func (s *Service) Start(c *context.Context, f *BaseForm) (interface{}, error) {
	visited := make(map[string]bool)
	var swap []string

	spider(f.BaseURL, visited, 3, &swap)
	fmt.Println("total links = ", len(swap))

	return removeDuplicateLink(swap), nil
}
