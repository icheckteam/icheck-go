package user

import (
	"github.com/icheckteam/icheck-go"
)

// Client is used to invoke /charges APIs.
type Client struct {
	B   icheck.Backend
	Key string
}

// Get returns the details of a user.
func Get(id string) (*icheck.User, error) {
	return getC().Get(id)
}

// Get ....
func (c Client) Get(id string) (*icheck.User, error) {
	var body *icheck.RequestValues
	var commonParams *icheck.Params
	user := &icheck.User{}
	err := c.B.Call("GET", "/users/"+id, c.Key, body, commonParams, user)
	return user, err
}

// List returns a list of users.
func List(params *icheck.UserListParams) (*icheck.UserList, error) {
	return getC().List(params)
}

// List ...
func (c Client) List(params *icheck.UserListParams) (*icheck.UserList, error) {
	var body *icheck.RequestValues
	var p *icheck.Params

	if params != nil {
		body = &icheck.RequestValues{}
		params.AppendTo(body)
		p = params.ToParams()
	}

	list := &icheck.UserList{}
	err := c.B.Call("GET", "/users", c.Key, body, p, list)
	return list, err
}

func getC() Client {
	return Client{icheck.GetBackend(icheck.APIBackend), icheck.Key}
}
