package address

import (
	"fmt"
	"strconv"

	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /account APIs.
type Client struct {
	B icheck.Backend
}

// List list all addresses
func (c *Client) List(params *icheck.Params) (*icheck.AddressListResp, error) {
	resp := &icheck.AddressListResp{}
	err := c.B.Call("GET", "/addresses", nil, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get get address detail
func (c *Client) Get(id string, params *icheck.Params) (*icheck.AddressResp, error) {
	resp := &icheck.AddressResp{}
	err := c.B.Call("GET", fmt.Sprintf("/addresses/%v", id), nil, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Create create an address
func (c *Client) Create(conf *icheck.AddressBody, params *icheck.Params) (*icheck.AddressResp, error) {
	body := &icheck.RequestValues{}
	if conf.Address != "" {
		body.Add("address", conf.Address)
	}
	if conf.City != 0 {
		body.Add("city", strconv.FormatInt(conf.City, 10))
	}
	if conf.District != 0 {
		body.Add("district", strconv.FormatInt(conf.District, 10))
	}
	if conf.Email != "" {
		body.Add("email", conf.Email)
	}
	resp := &icheck.AddressResp{}
	err := c.B.Call("POST", fmt.Sprintf("/addresses"), body, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Update update an address
func (c *Client) Update(id string, conf *icheck.AddressBody, params *icheck.Params) (*icheck.AddressResp, error) {
	body := &icheck.RequestValues{}
	if conf.Address != "" {
		body.Add("address", conf.Address)
	}
	if conf.City != 0 {
		body.Add("city", strconv.FormatInt(conf.City, 10))
	}
	if conf.District != 0 {
		body.Add("district", strconv.FormatInt(conf.District, 10))
	}
	if conf.Email != "" {
		body.Add("email", conf.Email)
	}
	resp := &icheck.AddressResp{}
	err := c.B.Call("PUT", fmt.Sprintf("/addresses/%v", id), body, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
