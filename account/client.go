package account

import (
	"fmt"
	"strconv"

	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /account APIs.
type Client struct {
	B icheck.Backend
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

// Login login user
func (c *Client) Logout(params *icheck.Params) (interface{}, error) {
	resp := make(map[string]interface{})
	err := c.B.Call("POST", "/logout", nil, params, resp)
	if err != nil {
		return nil, err
	}
	return resp["data"], nil
}

// LoginWithSocial ....
func (c *Client) LoginWithSocial(params *icheck.LoginSocialParams) (*icheck.AccessToken, error) {
	body := &icheck.RequestValues{}
	if params.Code != "" {
		body.Add("code", params.Code)
	}
	if params.TTL > 0 {
		body.Add("ttl", strconv.FormatInt(params.TTL, 10))
	}

	if params.Provider == "" {
		params.Provider = "facebook"
	}

	resp := &icheck.LoginResponse{}

	err := c.B.Call("GET", fmt.Sprintf("/auth/%s", params.Provider), body, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Register register an user
func (c *Client) Register(params *icheck.RegisterParams) (*icheck.UserResponse, error) {
	body := &icheck.RequestValues{}
	if params.Username != "" {
		body.Add("username", params.Username)
	}
	if params.Password != "" {
		body.Add("password", params.Password)
	}
	if params.Name != "" {
		body.Add("name", params.Name)
	}
	resp := &icheck.UserResponse{}
	err := c.B.Call("POST", "/register", body, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
