package icheck

// LoginResponse
type SearchResponse struct {
	Status int                    `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type SearchParams struct {
	Query string
	Type  string
	Limit int
	Skip  int
}
