package client

import (
	icheck "github.com/icheckteam/icheck-go"
	"github.com/icheckteam/icheck-go/account"
)

// API is the Icheck client. It contains all the different resources available.
type API struct {
	Account *account.Client
}

// Init initializes the Icheck client with the appropriate secret key
// as well as providing the ability to override the backend as needed.
func (a *API) Init(backend icheck.Backend) {
	if backend == nil {
		backend = icheck.GetBackend()
	}
	a.Account = &account.Client{B: backend}
}
