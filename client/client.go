// Package client provides a Icheck client for invoking APIs across all resources
package client

import (
	"github.com/icheckteam/icheck-go"
	"github.com/icheckteam/icheck-go/user"
)

// API is the Icheck client. It contains all the different resources available.
type API struct {
	User *user.Client
}

// Init initializes the Stripe client with the appropriate secret key
// as well as providing the ability to override the backend as needed.
func (a *API) Init(key string, backends *icheck.Backends) {
	a.User = &user.Client{B: backends.API}
}
