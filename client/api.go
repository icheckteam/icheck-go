package client

import (
	icheck "github.com/icheckteam/icheck-go"
	"github.com/icheckteam/icheck-go/account"
	"github.com/icheckteam/icheck-go/location"
	"github.com/icheckteam/icheck-go/search"
	"github.com/icheckteam/icheck-go/user"
)

// API is the Icheck client. It contains all the different resources available.
type API struct {
	Account  *account.Client
	User     *user.Client
	Search   *search.Client
	Location *location.Client
}

// Init initializes the Icheck client with the appropriate secret key
// as well as providing the ability to override the backend as needed.
func (a *API) Init(backend icheck.Backend) {
	if backend == nil {
		backend = icheck.GetBackend()
	}

	a.Account = &account.Client{B: backend}
	a.User = &user.Client{B: backend}
	a.Search = &search.Client{B: backend}
	a.Location = &location.Client{B: backend}
}

// New Api .....
func New() *API {
	api := &API{}
	api.Init(nil)
	return api
}
