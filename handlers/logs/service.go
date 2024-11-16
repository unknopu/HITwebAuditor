package logs

import (
	"auditor/entities"
	"auditor/handlers/common"
	"auditor/response"
)

// ServiceInterface service interface
type ServiceInterface interface {
	Insert(e *entities.Log)
	GetAllWithAdmin(f *GetAllWithAdminForm) (*response.AdminPage, error)
}

// Service user repo
type Service struct {
	rp RepoInterface
}

// NewService new service
func NewService() *Service {
	return &Service{
		rp: NewRepo(),
	}
}

// Insert insert
func (s *Service) Insert(e *entities.Log) {
	/*go func() {
		err := s.rp.Create(e)
		if err != nil {
			log.Error(err.Error())
		}
	}()*/
}

// GetAllWithAdminForm get all with admin form
type GetAllWithAdminForm struct {
	ID string `path:"id"`
	common.AdminPageQuery
}

// GetAllWithAdmin get all with admin
func (s *Service) GetAllWithAdmin(f *GetAllWithAdminForm) (*response.AdminPage, error) {
	us := []*entities.Log{}
	c, err := s.rp.FindAllWithAdminParameters(f, &us)
	if err != nil {
		return nil, err
	}
	return &response.AdminPage{
		PageInformation: &response.AdminPageInformation{
			Page:         f.GetPage(),
			Size:         f.PageSize(),
			NumberOfPage: int(c),
		},
		Entities: us,
	}, nil
}
