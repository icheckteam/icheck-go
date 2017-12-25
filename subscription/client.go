package account

import (
	icheck "github.com/icheckteam/icheck-go"
)

// Client is used to invoke /account APIs.
type Client struct {
	B icheck.Backend
}
