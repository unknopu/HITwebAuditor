package static_analysis

import (
	"auditor/app"
	"fmt"
	"strings"
	"sync"
	"text/scanner"

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
	content, err := fileContent(c)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func lexicalAnalysis(phpCode string) {
	var s scanner.Scanner
	s.Init(strings.NewReader(phpCode))
	s.Filename = "example.php"

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("Token: %s\tValue: %s\tPosition: %s\n", s.TokenText(), scanner.TokenString(tok), s.Pos())
	}
}
