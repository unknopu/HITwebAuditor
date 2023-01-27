package response

// AdminPageInformation admin page information
type AdminPageInformation struct {
	Page          int `json:"page,omitempty"`
	NumberOfPage  int `json:"number_of_page,omitempty"`
	Size          int `json:"size,omitempty"`
	TotalEntities int `json:"total_entities,omitempty"`
}

// AdminPage admin page
type AdminPage struct {
	PageInformation *AdminPageInformation `json:"page_information,omitempty"`
	Entities        interface{}           `json:"entities,omitempty"`
}
