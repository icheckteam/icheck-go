package search

import (
	"net/url"

	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /account APIs.
type Client struct {
	B icheck.Backend
}

// Search
func (c *Client) Search(params url.Values) (*icheck.SearchResponse, error) {
	body := &icheck.RequestValues{}
	if params.Get("type") != "" {
		body.Add("type", params.Get("type"))
	}
	if params.Get("query") != "" {
		body.Add("password", params.Get("query"))
	}
	if params.Get("limit") != "" {
		body.Add("limit", params.Get("limit"))
	}
	if params.Get("skip") != "" {
		body.Add("skip", params.Get("skip"))
	}
	resp := &icheck.SearchResponse{}
	err := c.B.Call("GET", "/search", body, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
