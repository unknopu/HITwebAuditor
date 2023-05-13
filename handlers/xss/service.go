package xss

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/core/utils"
	"auditor/entities"
	"log"
	"strings"
	"sync"
)

var (
	wg = sync.WaitGroup{}
	m  = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	Init(c *context.Context, f *XSSForm) (interface{}, error)
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

func (s *Service) Init(c *context.Context, f *XSSForm) (interface{}, error) {
	option := f.URLOptions()

	var reports []*entities.XSSReport
	testMsg := utils.RandomString(10)
	bodyTag := fetchTagBody(*option, testMsg)

	if !strings.ContainsAny(bodyTag, testMsg) {
		log.Println("test msg failure: ", testMsg)
		return reports, nil
	}

	var p []string
	var payloadSpliter string
	for _, payload := range payloads {
		bodyTag := fetchTagBody(*option, payload)
		if strings.ContainsAny(bodyTag, payload) {
			if strings.Contains(payload, payloadSpliter) && payloadSpliter != "" {
				continue
			}

			p = append(p, payload)
			payloadSpliter = payload[0:4]
		}
		if len(p) > 3 {
			break
		}
	}

	if len(p) > 0 {
		reports = append(reports, &entities.XSSReport{
			Location:       f.URL,
			Payload:        p,
			Level:          entities.HIGH,
			Type:           entities.Injection,
			Vaulnerability: entities.CrossSiteScripting,
		})
	}

	return buildPageInfomation(reports), nil
}
