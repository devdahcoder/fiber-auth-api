package types

type Metadata struct {
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	PerPage      int `json:"per_page"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
	TotalCurrentRecords int `json:"total_current_records"`
}

type QueryTypeValidationFunction func(string) bool