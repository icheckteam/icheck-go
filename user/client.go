package user

import (
	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /users APIs.
type Client struct {
	B icheck.Backend
}

// Login login user
func (c *Client) Get(userID string, params *icheck.Params) (*icheck.AccessToken, error) {
	resp := &icheck.LoginResponse{}
	err := c.B.Call("GET", "/users/"+userID, nil, params, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// List ...
func (c *Client) List(params *icheck.UserListParams) ([]icheck.User, error) {
	body := &icheck.RequestValues{}

	if len(params.IcheckID) > 0 {
		for _, userID := range params.IcheckID {
			body.Add("icheck_id", userID)
		}
	}

	resp := &icheck.UserListResponse{}
	err := c.B.Call("GET", "/users", body, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}
