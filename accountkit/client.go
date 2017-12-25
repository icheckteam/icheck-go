package accountkit

import (
	"strconv"

	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /account APIs.
type Client struct {
	B icheck.Backend
}

func (c *Client) Login(params *icheck.AccountKitLoginParams) (*icheck.AccessToken, error) {
	body := &icheck.RequestValues{}

	if params.Code != "" {
		body.Add("code", params.Code)
	}

	if params.Name != "" {
		body.Add("name", params.Name)
	}

	if params.Password != "" {
		body.Add("password", params.Password)
	}

	if params.TTL > 0 {
		body.Add("ttl", strconv.FormatInt(params.TTL, 10))
	}

	resp := &icheck.LoginResponse{}

	err := c.B.Call("POST", "/accountkit/login", body, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) ResetPassword(params *icheck.AccountKitResetPasswordParams) (*icheck.AccountKitResetPasswordResponse, error) {
	body := &icheck.RequestValues{}

	if params.Code != "" {
		body.Add("code", params.Code)
	}

	if params.Password != "" {
		body.Add("password", params.Password)
	}

	resp := &icheck.AccountKitResetPasswordResponse{}

	err := c.B.Call("POST", "/accountkit/login", body, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
