package account

import (
	"strconv"

	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /account APIs.
type Client struct {
	B   icheck.Backend
	Key string
}

// Me get current user
func (c *Client) Me(params *icheck.Params) (*icheck.User, error) {
	res := &icheck.UserResponse{}
	err := c.B.Call("GET", "/account", nil, params, res)

	if err != nil {
		return nil, err
	}
	return res.User, nil
}

// Login login user
func (c *Client) Login(params *icheck.LoginParams) (*icheck.AccessToken, error) {
	body := &icheck.RequestValues{}
	if params.Username != "" {
		body.Add("username", params.Username)
	}
	if params.Password != "" {
		body.Add("password", params.Password)
	}
	if params.TTL > 0 {
		body.Add("ttl", strconv.FormatInt(params.TTL, 10))
	}
	resp := &icheck.LoginResponse{}
	err := c.B.Call("POST", "/login", body, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
