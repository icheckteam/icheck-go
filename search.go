package icheck

// LoginResponse
type SearchResponse struct {
	Status int
	Data   map[string]interface{}
}

type SearchParams struct {
	Query string
	Type  string
	Limit int
	Skip  int
}
