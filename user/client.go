package user

import (
	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /users APIs.
type Client struct {
	B icheck.Backend
}

// Login login user
func (c *Client) Get(userID string, params *icheck.Params) (*icheck.User, error) {
	resp := &icheck.UserResponse{}
	err := c.B.Call("GET", "/users/"+userID, nil, params, resp)
	if err != nil {
		return nil, err
	}
	return resp.User, nil
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

// Update ...
func (c *Client) Update(data *icheck.UserUpdateParams, params *icheck.Params) (*icheck.User, error) {
	body := &icheck.RequestValues{}
	if params.Name != "" {
		body.Add("name", params.Name)
	}
	if params.Avatar != "" {
		body.Add("avatar", params.Avatar)
	}

	if params.Cover != "" {
		body.Add("cover", params.Cover)
	}

	resp := &icheck.UserResponse{}
	err := c.B.Call("POST", "/account", body, params, resp)
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}
