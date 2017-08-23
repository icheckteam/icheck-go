package client

import (
	"log"
	"testing"

	icheck "github.com/icheckteam/icheck-go"
)

func TestAPI(t *testing.T) {
	api := New()

	user, err := api.Account.Me(&icheck.Params{
		AccessToken: "1",
	})

	if err != nil {
		t.Fatal(err)
	}

	log.Print(user)
}
