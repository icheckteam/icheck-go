package icheck

// LocationsResponse
type LocationsResponse struct {
	Status int                      `json:"status"`
	Data   []map[string]interface{} `json:"data"`
}

// LocationResponse
type LocationResponse struct {
	Status int                    `json:"status"`
	Data   map[string]interface{} `json:"data"`
}
