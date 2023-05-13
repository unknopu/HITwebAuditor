package static_analysis

import (
	"auditor/app"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c echo.Context, f *StaticAnalysisForm) (interface{}, error)
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

func (s *Service) Init(c echo.Context, f *StaticAnalysisForm) (interface{}, error) {

	content, err := fetchSourceCode(c, f)
	if err != nil {
		return nil, err
	}
	sourceCode := validatePHPSource(content)
	lexicalAnalysis(sourceCode)

	return sourceCode, nil
}
