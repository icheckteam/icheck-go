package location

import (
	"fmt"
	"net/url"

	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /account APIs.
type Client struct {
	B icheck.Backend
}

// Me get current user
func (c *Client) List(params url.Values) (*icheck.LocationsResponse, error) {
	body := &icheck.RequestValues{}

	if params.Get("parent") != "" {
		body.Add("parent", params.Get("parent"))
	}

	if params.Get("type") != "" {
		body.Add("type", params.Get("type"))
	} else {
		body.Add("type", "city")
	}
	resp := &icheck.LocationsResponse{}
	err := c.B.Call("GET", "/locations", body, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Me get current user
func (c *Client) Get(id string) (*icheck.LocationResponse, error) {
	resp := &icheck.LocationResponse{}
	err := c.B.Call("GET", fmt.Sprintf("/locations/%v", id), nil, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
